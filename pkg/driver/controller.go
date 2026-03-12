package driver

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type ControllerService struct {
}

func (cs *ControllerService) CreateVolume(*csi.CreateVolumeRequest) *csi.CreateVolumeResponse {
	fmt.Print("This is createvol")
	return nil
}

func (cs *ControllerService) DeleteVolume(*csi.DeleteVolumeRequest) *csi.DeleteVolumeResponse {
	fmt.Print("This is deletevol")
	return nil
}

func (cs *ControllerService) ControllerPublishVolume(*csi.ControllerPublishVolumeRequest) *csi.ControllerPublishVolumeResponse {
	fmt.Print("This is controllerpubvol")
	return nil
}

func (cs *ControllerService) ControllerUnpublishVolume(*csi.ControllerUnpublishVolumeRequest) *csi.ControllerUnpublishVolumeResponse {
	fmt.Print("This is controllerunpubvol")
	return nil
}
