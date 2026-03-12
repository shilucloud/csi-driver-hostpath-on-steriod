package driver

import "fmt"

type Driver struct {
	controller *ControllerService
	node       *NodeService
	options    Options
	srv        string
	healthy    bool
}

type Options struct {
	Mode     Mode
	Endpoint string
	Name     string
}

func NewDriver(options *Options) *Driver {
	driver := Driver{
		controller: &ControllerService{},
		node:       &NodeService{},
		options:    *options,
		healthy:    false,
	}
	return &driver
}

func (d *Driver) Run() {

	switch d.options.Mode {
	case ControllerMode:
		fmt.Print("controller")

	case NodeMode:
		fmt.Print("controller")

	default:
		fmt.Print("ALL")
	}

}

func (d *Driver) Stop() {
	fmt.Print("Stopping the server")
}
