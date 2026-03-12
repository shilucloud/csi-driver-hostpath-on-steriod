package driver

import "fmt"

type Driver struct {
	controller    *ControllerService
	node          *NodeService
	driverName    string
	driverVersion string
	options       Options
	srv           string
	healthy       bool
}

type Options struct {
	Mode          Mode
	Endpoint      string
	Name          string
	driverVersion string
}

func NewDriver(options *Options) (*Driver, error) {
	driver := Driver{
		driverName:    options.Name,
		driverVersion: options.driverVersion,
		options:       *options,
	}

	switch options.Mode {
	case ControllerMode:
		driver.controller = NewControllerService()
	case NodeMode:
		driver.node = NewNodeService()
	case AllMode:
		driver.controller = NewControllerService()
		driver.node = NewNodeService()
	default:
		return nil, fmt.Errorf("unknown mode: %s", options.Mode)
	}

	return &driver, nil
}

func (d *Driver) Run() error {

	switch d.options.Mode {
	case ControllerMode:
		fmt.Print("controller")

	case NodeMode:
		fmt.Print("controller")

	default:
		fmt.Print("ALL")
	}

	return nil

}

func (d *Driver) Stop() {
	fmt.Print("Stopping the server")
}
