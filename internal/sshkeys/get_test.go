package sshkeys_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"slices"
	"testing"
)

func getRespMock(sshKeyID int, opts *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if sshKeyID != 1 {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key id does not match: want %d, got %d", 1, sshKeyID)
	}
	wantFields := []string{"ssh_key", "email"}
	if !slices.Equal(opts.Fields, wantFields) {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key get options do not match: want %v, got %v", wantFields, opts.Fields)
	}
	return stubs.SSHKeyWithEmail, nil, nil
}

func TestGetSSHKey(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.GetResp = getRespMock
	})

	f.RootCmd.SetArgs([]string{"ssh-key", "get", "1"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Label", "User", "Fingerprint", "Created"}
	wantTD := [][]string{{"1", "test_label", "test@example.com", "test_fingerprint", "test_created"}}

	f.CheckOutput(stubs.SSHKeyWithEmail, wantTH, wantTD)
}

func TestGetSSHKeyAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.GetResp = getRespMock
	})
	f.RootCmd.SetArgs([]string{"ssh-key", "get", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get ssh-key: ssh key id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
