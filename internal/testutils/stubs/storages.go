package stubs

import "github.com/cherryservers/cherrygo/v3"

var BlockStorage = cherrygo.BlockStorage{
	ID:   1,
	Name: "cs-volume-1-1",
	Href: "/storage/1",
	Region: cherrygo.Region{
		ID:         1,
		Name:       "test",
		Slug:       "test",
		RegionIso2: "LT",
		BGP: cherrygo.RegionBGP{
			Hosts: []string{"1.1.1.1", "2.2.2.2"},
			Asn:   12345,
		},
		Location: "Lithuania, Šiauliai",
	},
	Size:          1,
	AllowEditSize: true,
	Unit:          "GB",
	Description:   "test",
	AttachedTo: cherrygo.AttachedTo{
		ID:       1,
		Hostname: "test",
		Href:     "/servers/1",
	},
}
