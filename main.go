package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/jdextraze/docker-machine-driver-atlanticnet/driver"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver("", ""))
}
