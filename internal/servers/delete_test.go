package servers_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func deleteRespMock(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	if serverID != 1 {
		return cherrygo.Server{}, nil,
			fmt.Errorf("server id does not match: want %d, got %d", 1, serverID)
	}
	return cherrygo.Server{}, nil, nil
}

func TestDeleteServerWithForce(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeServersService.DeleteResp = deleteRespMock
	})
	f.RootCmd.SetArgs([]string{"server", "delete", "1", "-f"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	f.CheckOutputLine(fmt.Sprintln("Server", 1, "successfully deleted."))
}

func TestDeleteServerWithForceAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeServersService.DeleteResp = deleteRespMock
	})
	f.RootCmd.SetArgs([]string{"server", "delete", "2", "--force"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not delete Server: server id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)

}
func TestDeleteServerWithInvalidIDError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {})
	f.RootCmd.SetArgs([]string{"server", "delete", "wrong", "--force"})
	err := f.RootCmd.Execute()
	const wantMsg = "invalid server ID: strconv.Atoi: parsing \"wrong\": invalid syntax"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
