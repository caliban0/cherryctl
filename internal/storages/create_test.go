package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func TestCreateStorageKeyMissingFlagsError(t *testing.T) {
	for _, v := range [][]string{
		{"storage", "create", "--project-id", "1", "--size", "1"},
		{"storage", "create", "--region", "eu_nord_1", "--project-id", "1"},
		{"storage", "create", "--region", "eu_nord_1", "--size", "1"},
	} {
		testutils.VerifyMissingFlagsError(t, v)
	}
}

func createRespMock(request *cherrygo.CreateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	wantReq := cherrygo.CreateStorage{
		ProjectID:   1,
		Description: "test",
		Size:        1,
		Region:      "test",
	}
	if request == nil || wantReq != *request {
		return cherrygo.BlockStorage{}, nil,
			fmt.Errorf("bad creation request, want: %v, got %v", &wantReq, request)
	}
	return stubs.BlockStorage, nil, nil
}

func TestCreateStorage(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "create", "--project-id", "1", "--size", "1",
		"--region", "test", "--description", "test"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Size", "Region", "Description"}
	wantTD := [][]string{{"1", "1 GB", "test", "test"}}
	f.CheckOutput(stubs.BlockStorage, wantTH, wantTD)
}

func TestCreateStorageAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"storage", "create", "--project-id", "-1", "--size", "-1",
		"--region", "test-error", "--description", "test-error"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not create storage: bad creation request, want: &{1 test 1 test}, got &{-1 test-error -1 test-error}"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
