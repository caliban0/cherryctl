package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"slices"
	"testing"
)

func getResp(storageID int, opts *cherrygo.GetOptions) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	if storageID != 1 {
		return cherrygo.BlockStorage{}, nil, fmt.Errorf("storage id does not match: want %d, got %d", 1, storageID)
	}
	wFields := []string{"storage", "region", "id", "hostname"}
	if !slices.Equal(opts.Fields, wFields) {
		return cherrygo.BlockStorage{}, nil,
			fmt.Errorf("storage get options do not match: want %v, got %v", wFields, opts.Fields)
	}
	return stubs.BlockStorage, nil, nil
}

func TestGetStorage(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.GetResp = getResp
	})
	f.RootCmd.SetArgs([]string{"storage", "get", "1"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	wantTH := []string{"ID", "Size", "Region", "Description", "Attached to"}
	wantTD := [][]string{{"1", "1 GB", "test", "test", "test"}}
	f.CheckOutput(stubs.BlockStorage, wantTH, wantTD)
}

func TestGetStorageAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.GetResp = getResp
	})
	f.RootCmd.SetArgs([]string{"storage", "get", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get storage: storage id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
