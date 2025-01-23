package storages_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"slices"
	"testing"
)

// Expect request with StorageID == 1, AttachTo == 1.
func attachRespMock(request *cherrygo.AttachTo) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	wantReq := cherrygo.AttachTo{StorageID: 1, AttachTo: 1}
	if request == nil || wantReq != *request {
		return cherrygo.BlockStorage{}, nil,
			fmt.Errorf("bad creation request, want: %v, got %v", &wantReq, request)
	}
	return stubs.BlockStorage, nil, nil
}

// Expect projectID == 1, with Fields opts {"id", "name", "hostname"}.
func listServersMock(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.Server, *cherrygo.Response, error) {
	if projectID != 1 {
		return []cherrygo.Server{}, nil,
			fmt.Errorf("project id does not match: want %d, got %d", 1, projectID)
	}
	wantFields := []string{"id", "name", "hostname"}
	if !slices.Equal(opts.Fields, wantFields) {
		return []cherrygo.Server{}, nil,
			fmt.Errorf("server list options do not match: want %v, got %v", wantFields, opts.Fields)
	}
	return []cherrygo.Server{stubs.ServerWithIDAndNames}, nil, nil
}

// Known issue: can't attach by server ID flag.
func TestAttachStorage(t *testing.T) {
	wHostname := stubs.BlockStorage.AttachedTo.Hostname
	for _, v := range [][]string{{"storage", "attach", "1", "--server-id", "1"},
		{"storage", "attach", "1", "--project-id", "1", "--server-hostname", wHostname}} {
		f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
			fakeServicer.FakeStorageService.AttachResp = attachRespMock
			fakeServicer.FakeServersService.ListResp = listServersMock
		})
		f.RootCmd.SetArgs(v)
		if err := f.RootCmd.Execute(); err != nil {
			t.Error(err)
			return
		}
		f.CheckOutputLine(fmt.Sprintln("Storage", 1, "successfully attached to", wHostname))
	}
}

// Same issue as TestAttachStorage.
func TestAttachStorageAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeStorageService.AttachResp = attachRespMock
		fakeServicer.FakeServersService.ListResp = listServersMock
	})
	f.RootCmd.SetArgs([]string{"storage", "attach", "1", "--server-id", "2"})
	err := f.RootCmd.Execute()
	testutils.VerifyErrorMsg(t, err, "Could not atach storage to a server: "+
		"bad creation request, want: &{1 1}, got &{1 2}")
}

func TestAttachStorageWithNoIdOrHostnameError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {})
	f.RootCmd.SetArgs([]string{"storage", "attach", "1"})
	err := f.RootCmd.Execute()
	testutils.VerifyErrorMsg(t, err, "either server-id or server-hostname should be set")
}

func TestAttachStorageByHostnameNoProjectIdError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {})
	f.RootCmd.SetArgs([]string{"storage", "attach", "1", "--server-hostname", "test"})
	err := f.RootCmd.Execute()
	testutils.VerifyErrorMsg(t, err, "when using --server-hostname, project-id is a required argument")
}

// Known issue: attach request proceeds, when a server with the provided hostname can't be found
// (for a valid project ID). This is a bug in utils.ServerHostnameToID
func TestAttachStorageServerByHostnameNotFoundError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeServersService.ListResp = listServersMock
	})
	f.RootCmd.SetArgs([]string{"storage", "attach", "1", "--server-hostname", "missing", "--project-id", "1"})
	err := f.RootCmd.Execute()
	testutils.VerifyErrorMsg(t, err, "Could not get a Server")
}
