package stubs

import "github.com/cherryservers/cherrygo/v3"

var ServerWithIDAndNames = cherrygo.Server{
	ID:       1,
	Name:     "test",
	Hostname: "test",
}

var Server = cherrygo.Server{
	ID:       1,
	Name:     "test_name",
	Hostname: "test_hostname",
	Image:    "ubuntu_22_04",
	State:    "deployed",
	Region:   cherrygo.Region{Name: "eu_nord_1"},
}
