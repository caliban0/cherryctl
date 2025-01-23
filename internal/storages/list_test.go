package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"slices"
	"testing"
)

func listRespMock(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.BlockStorage, *cherrygo.Response, error) {
	if projectID != 1 {
		return []cherrygo.BlockStorage{}, nil, fmt.Errorf("project id does not match: want %d, got %d", 1, projectID)
	}
	wantFields := []string{"storage", "region", "id", "hostname"}
	if !slices.Equal(opts.Fields, wantFields) {
		return []cherrygo.BlockStorage{}, nil,
			fmt.Errorf("storage get options do not match: want %v, got %v", wantFields, opts.Fields)
	}
	return []cherrygo.BlockStorage{stubs.BlockStorage}, nil, nil
}

func TestListStorages(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.ListResp = listRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "list", "-p", "1"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Size", "Region", "Description", "Attached to"}
	wantTD := [][]string{{"1", "1 GB", "test", "test", "test"}}
	f.CheckOutput([]cherrygo.BlockStorage{stubs.BlockStorage}, wantTH, wantTD)
}

func TestListStoragesAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.ListResp = listRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "list", "-p", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get storage list: project id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}

func TestListStoragesMissingRequiredFlags(t *testing.T) {
	testutils.VerifyMissingFlagsError(t, []string{"storage", "list"})
}
