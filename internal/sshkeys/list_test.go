package sshkeys_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"slices"
	"strconv"
	"testing"
)

func projectListRespMock(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	if projectID != 1 {
		return []cherrygo.SSHKey{}, nil,
			fmt.Errorf("project id does not match: want %d, got %d", 1, projectID)
	}
	wantFields := []string{"ssh_key", "email"}
	if !slices.Equal(opts.Fields, wantFields) {
		return []cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key get options do not match: want %v, got %v", wantFields, opts.Fields)
	}
	return []cherrygo.SSHKey{stubs.SSHKeyWithEmail}, nil, nil
}

func listRespMock(opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	wantFields := []string{"ssh_key", "email"}
	if !slices.Equal(opts.Fields, wantFields) {
		return []cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key get options do not match: want %v, got %v", wantFields, opts.Fields)
	}
	return []cherrygo.SSHKey{stubs.SSHKeyWithEmail}, nil, nil
}

func TestListSSHKeys(t *testing.T) {
	wantR := stubs.SSHKeyWithEmail
	for _, v := range [][]string{{"ssh-key", "list"},
		{"ssh-key", "list", "--project-id", "1"},
		{"ssh-key", "list", "-p", "1"}} {
		f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
			fakeServicer.FakeSSHKeysService.ListResp = listRespMock
			fakeServicer.FakeProjectService.ListSSHKeysResp = projectListRespMock
		})

		f.RootCmd.SetArgs(v)
		if err := f.RootCmd.Execute(); err != nil {
			t.Error(err)
			return
		}
		wantTH := []string{"ID", "Label", "User", "Fingerprint", "Created"}
		wantTD := [][]string{{strconv.Itoa(wantR.ID), wantR.Label,
			wantR.User.Email, wantR.Fingerprint, wantR.Created}}

		f.CheckOutput([]cherrygo.SSHKey{wantR}, wantTH, wantTD)
	}

}

func TestListSSHKeysAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeProjectService.ListSSHKeysResp = projectListRespMock
	})
	f.RootCmd.SetArgs([]string{"ssh-key", "list", "-p", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get ssh-keys list: project id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
