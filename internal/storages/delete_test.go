package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func deleteRespMock(storageID int) (*cherrygo.Response, error) {
	if storageID != 1 {
		return nil,
			fmt.Errorf("storage id does not match: want %d, got %d", 1, storageID)
	}
	return nil, nil
}

func TestDeleteStorageWithForce(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.DeleteResp = deleteRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "delete", "1", "-f"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	f.CheckOutputLine(fmt.Sprintln("Storage", 1, "successfully deleted."))
}

func TestDeleteStorageWithForceAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.DeleteResp = deleteRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "delete", "2", "--force"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not delete storage: storage id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
