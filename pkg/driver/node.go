package driver

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/shilucloud/csi-driver-hostpath-on-steriod/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

var (
	nodeCaps = []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
		csi.NodeServiceCapability_RPC_VOLUME_MOUNT_GROUP,
	}
)

type NodeService struct {
	csi.UnimplementedNodeServer
}

func (ns *NodeService) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.InfoS("Received NodeStageVolume request", "volID", req.VolumeId, "stagingPath", req.StagingTargetPath)

	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "volume ID must be provided")
	}
	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "staging target path must be provided")
	}
	if req.PublishContext["imgPath"] == "" {
		return nil, status.Error(codes.InvalidArgument, "imgPath must be provided in publish context")
	}
	if req.VolumeContext["fsType"] == "" {
		return nil, status.Error(codes.InvalidArgument, "fsType must be provided in volume context")
	}
	if req.VolumeContext["byteSize"] == "" {
		return nil, status.Error(codes.InvalidArgument, "byteSize must be provided in volume context")
	}

	byteSize, err := util.StrToInt(req.VolumeContext["byteSize"])
	if err != nil {
		klog.ErrorS(err, "failed to parse byteSize", "byteSize", req.VolumeContext["byteSize"])
		return nil, status.Errorf(codes.Internal, "failed to parse byteSize: %v", err)
	}

	imgPath := req.PublishContext["imgPath"]
	fsType := req.VolumeContext["fsType"]
	stagingPath := req.StagingTargetPath

	// 1. Create the image file
	if err := util.CreateImageFile(imgPath, byteSize); err != nil {
		klog.ErrorS(err, "failed to create image file", "imgPath", imgPath)
		return nil, status.Errorf(codes.Internal, "failed to create image file: %v", err)
	}
	klog.InfoS("Image file ready", "imgPath", imgPath)

	// 2. Attach loop device
	devicePath, err := util.AttachLoopDevice(imgPath)
	if err != nil {
		klog.ErrorS(err, "failed to attach loop device", "imgPath", imgPath)
		return nil, status.Errorf(codes.Internal, "failed to attach loop device: %v", err)
	}
	klog.InfoS("Loop device attached", "devicePath", devicePath)

	// 3. Format filesystem
	if err := util.MakeFs(devicePath, fsType); err != nil {
		klog.ErrorS(err, "failed to make filesystem", "devicePath", devicePath, "fsType", fsType)
		return nil, status.Errorf(codes.Internal, "failed to make filesystem: %v", err)
	}
	klog.InfoS("Filesystem created", "devicePath", devicePath, "fsType", fsType)

	// 4. Mount to staging path
	if err := util.Mount(devicePath, stagingPath, fsType); err != nil {
		klog.ErrorS(err, "failed to mount to staging path", "devicePath", devicePath, "stagingPath", stagingPath)
		return nil, status.Errorf(codes.Internal, "failed to mount: %v", err)
	}
	klog.InfoS("Mounted to staging path", "devicePath", devicePath, "stagingPath", stagingPath)

	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *NodeService) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.InfoS("Received NodeUnstageVolume request", "volID", req.VolumeId, "stagingPath", req.StagingTargetPath)

	// 1. detach loop device first (while still mounted so findmnt can find it)
	if err := util.DetachLoopDevice(req.StagingTargetPath); err != nil {
		klog.ErrorS(err, "failed to detach loop device", "stagingPath", req.StagingTargetPath)
		// non-fatal, continue
	}
	klog.InfoS("Loop device detached", "stagingPath", req.StagingTargetPath)

	// 2. then unmount
	if err := util.UnmountOnly(req.StagingTargetPath); err != nil {
		klog.ErrorS(err, "failed to unmount staging path", "stagingPath", req.StagingTargetPath)
		return nil, status.Errorf(codes.Internal, "failed to unmount: %v", err)
	}
	klog.InfoS("NodeUnstageVolume complete", "volID", req.VolumeId)

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *NodeService) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.InfoS("Received NodePublishVolume request", "volID", req.VolumeId, "stagingPath", req.TargetPath)

	err := util.BindMount(req.StagingTargetPath, req.TargetPath)
	if err != nil {
		klog.ErrorS(err, "failed to mount to target path", "stagingPath", req.StagingTargetPath, "targetPath", req.TargetPath)
		return nil, status.Errorf(codes.Internal, "failed to mount: %v", err)
	}

	klog.InfoS("Mounted to Target path", "stagingPath", req.StagingTargetPath, "targetPath", req.TargetPath)

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *NodeService) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.InfoS("Received NodeUnpublishVolume request", "volID", req.VolumeId, "TargetPath", req.TargetPath)
	err := util.Unmount(req.TargetPath)
	if err != nil {
		klog.ErrorS(err, "failed to unmount target path", "targetPath", req.TargetPath)
		return nil, status.Errorf(codes.Internal, "failed to unmount: %v", err)
	}
	klog.InfoS("Unmounted from Target path", "targetPath", req.TargetPath)

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (n *NodeService) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return &csi.NodeExpandVolumeResponse{}, nil
}

func (n *NodeService) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {

	caps := make([]*csi.NodeServiceCapability, 0, len(nodeCaps))
	for _, capability := range nodeCaps {
		c := &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: capability,
				},
			},
		}
		caps = append(caps, c)
	}
	return &csi.NodeGetCapabilitiesResponse{Capabilities: caps}, nil
}

func (n *NodeService) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		return nil, status.Error(codes.Internal, "NODE_NAME env var not set")
	}

	return &csi.NodeGetInfoResponse{
		NodeId:            nodeName,
		MaxVolumesPerNode: 5,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"kubernetes.io/hostname": nodeName,
			},
		},
	}, nil
}

func (n *NodeService) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{}, nil
}

func NewNodeService() *NodeService {
	return &NodeService{}
}
