package main

import (
	"fmt"

	"github.com/shilucloud/csi-driver-hostpath-on-steriod/pkg/driver"
)

func main() {
	fmt.Println("This is a new go module")
	d := driver.NewDriver(&driver.Options{
		Mode:     "controller",
		Endpoint: "e",
		Name:     "Name",
	})
	d.Run()

}
