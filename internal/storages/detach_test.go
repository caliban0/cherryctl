package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func detachRespMock(storageID int) (*cherrygo.Response, error) {
	if storageID != 1 {
		return nil,
			fmt.Errorf("storage id does not match: want %d, got %d", 1, storageID)
	}
	return nil, nil
}

func TestDetachStorage(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.DetachResp = detachRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "detach", "1"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	f.CheckOutputLine(fmt.Sprintln("Storage volume", 1, "detached successfully."))
}

func TestDetachStorageAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.DetachResp = detachRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "detach", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not detach storage from a server: storage id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
