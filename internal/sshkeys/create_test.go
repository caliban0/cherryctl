package sshkeys_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func TestCreateSSHKeyMissingFlagsError(t *testing.T) {
	for _, v := range [][]string{
		{"ssh-key", "create", "--label", "test"},
		{"ssh-key", "create", "--key", "test"},
	} {
		testutils.VerifyMissingFlagsError(t, v)
	}
}

func createRespMock(request *cherrygo.CreateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	wantReq := cherrygo.CreateSSHKey{Label: "test_label", Key: "test_key"}
	if request == nil || *request != wantReq {
		return cherrygo.SSHKey{}, nil,
			fmt.Errorf("bad creation request, want: %v, got %v", &wantReq, request)
	}
	return stubs.SSHKey, nil, nil
}

func TestCreateSSHKey(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"ssh-key", "create", "--label", "test_label", "--key", "test_key"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Label", "Fingerprint", "Created"}
	wantTD := [][]string{{"1", "test_label", "test_fingerprint", "test_created"}}
	f.CheckOutput(stubs.SSHKey, wantTH, wantTD)
}

func TestCreateSSHKeyAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeSSHKeysService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"ssh-key", "create", "--label", "wrong", "--key", "wrong"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not create SSH key: bad creation request, want: &{test_label test_key}, got &{wrong wrong}"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
