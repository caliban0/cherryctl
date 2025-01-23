package stubs

import "github.com/cherryservers/cherrygo/v3"

var SSHKeyWithEmail = cherrygo.SSHKey{
	ID:          1,
	Label:       "test_label",
	Key:         "test_key",
	Fingerprint: "test_fingerprint",
	Updated:     "test_updated",
	Created:     "test_created",
	Href:        "test_href",
	User:        cherrygo.User{Email: "test@example.com"}}

var SSHKey = cherrygo.SSHKey{
	ID:          1,
	Label:       "test_label",
	Key:         "test_key",
	Fingerprint: "test_fingerprint",
	Updated:     "test_updated",
	Created:     "test_created",
	Href:        "test_href",
}
