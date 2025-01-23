package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func updateRespMock(request *cherrygo.UpdateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	wantReq := cherrygo.UpdateStorage{
		StorageID:   1,
		Description: "test",
		Size:        1,
	}
	if request == nil || wantReq != *request {
		return cherrygo.BlockStorage{}, nil,
			fmt.Errorf("bad update request, want: %v, got %v", &wantReq, request)
	}
	return stubs.BlockStorage, nil, nil
}

func TestUpdateStorage(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.UpdateResp = updateRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "update", "1", "--size", "1", "--description", "test"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Size", "Region", "Description", "Attached to"}
	wantTD := [][]string{{"1", "1 GB", "test", "test", "test"}}
	f.CheckOutput(stubs.BlockStorage, wantTH, wantTD)
}

func TestUpdateStorageAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.UpdateResp = updateRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "update", "2", "--size", "2", "--description", "wrong"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not update storage: bad update request, want: &{1 1 test}, got &{2 2 wrong}"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
