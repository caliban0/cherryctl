package sshkeys_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"strconv"
	"testing"
)

func updateRespMock(sshKeyID int, request *cherrygo.UpdateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	wantR := stubs.SSHKey
	wantReq := cherrygo.UpdateSSHKey{Label: &wantR.Label, Key: &wantR.Key}
	if sshKeyID != wantR.ID {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key id does not match: want %d, got %d", wantR.ID, sshKeyID)
	}

	if request == nil || *request.Key != wantR.Key || *request.Label != wantR.Label {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("bad update request, want: %v, got %v", wantReq, request)
	}
	return stubs.SSHKey, nil, nil
}

func TestUpdateSSHKey(t *testing.T) {
	wantR := stubs.SSHKey
	id := strconv.Itoa(wantR.ID)
	for _, v := range [][]string{
		{"ssh-key", "update", id, "--ssh-key-id", id, "--label", wantR.Label, "--key", wantR.Key},
		{"ssh-key", "update", id, "-i", id, "--label", wantR.Label, "--key", wantR.Key},
	} {
		f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
			fakeServicer.FakeSSHKeysService.UpdateResp = updateRespMock
		})
		f.RootCmd.SetArgs(v)

		if err := f.RootCmd.Execute(); err != nil {
			t.Error(err)
			return
		}
		wantTH := []string{"ID", "Label", "Fingerprint", "Created"}
		wantTD := [][]string{{id, wantR.Label,
			wantR.Fingerprint, wantR.Created}}

		f.CheckOutput(wantR, wantTH, wantTD)

	}
}

func TestUpdateSSHKeyAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.UpdateResp = updateRespMock
	})
	id := strconv.Itoa(stubs.SSHKey.ID)
	f.RootCmd.SetArgs([]string{"ssh-key", "update", id, "--ssh-key-id", id, "--label", "wrong", "--key", "wrong"})
	err := f.RootCmd.Execute()

	// Can't check the exact message, because update request fields are pointers.
	const wantMsg = "Could not update SSH key: bad update request, want: (.*)"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
