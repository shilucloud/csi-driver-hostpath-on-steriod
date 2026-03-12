package driver

type Mode string

const (
	// ControllerMode is the mode that only starts the controller service.
	ControllerMode Mode = "controller"
	// NodeMode is the mode that only starts the node service.
	NodeMode Mode = "node"
	// AllMode is the mode that only starts both the controller and the node service.
	AllMode Mode = "all"

	// MetadataLabelerMode is the mode that starts the metadata labeler.
	MetadataLabelerMode = "metadataLabeler"
)
