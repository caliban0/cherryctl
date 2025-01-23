package sshkeys_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func TestDeleteSSHKeyWithForceMissingArgError(t *testing.T) {
	testutils.VerifyMissingFlagsError(t, []string{"ssh-key", "delete"})
}

func deleteRespMock(sshKeyID int) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if sshKeyID != 1 {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("ssh key id does not match: want %d, got %d", 1, sshKeyID)
	}
	return cherrygo.SSHKey{}, nil, nil
}

func TestDeleteSSHKeyWithForce(t *testing.T) {
	for _, v := range [][]string{{"ssh-key", "delete", "-i", "1", "--force"},
		{"ssh-key", "delete", "--ssh-key-id", "1", "-f"},
		{"ssh-key", "delete", "-i", "1", "-f"}} {
		f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
			fakeServicer.FakeSSHKeysService.DeleteResp = deleteRespMock
		})
		f.RootCmd.SetArgs(v)

		if err := f.RootCmd.Execute(); err != nil {
			t.Error(err)
			return
		}
		f.CheckOutputLine(fmt.Sprintln("SSH key", 1, "successfully deleted."))
	}
}

func TestDeleteSSHKeyWithForceAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.DeleteResp = deleteRespMock
	})
	f.RootCmd.SetArgs([]string{"ssh-key", "delete", "--ssh-key-id", "2", "-f"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not delete SSH key: ssh key id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
