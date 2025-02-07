package servers_test

import (
	"encoding/base64"
	"maps"
	"slices"

	//"encoding/base64"
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"os"
	"testing"
)

func TestCreateStorageKeyMissingFlagsError(t *testing.T) {
	for _, v := range [][]string{
		{"server", "create", "--project-id", "1", "--plan", "cloud_vps_1"},
		{"server", "create", "--region", "eu_nord_1", "--project-id", "1"},
		{"server", "create", "--region", "eu_nord_1", "--plan", "cloud_vps_1"},
	} {
		testutils.VerifyMissingFlagsError(t, v)
	}
}

func equalCreateServerRequest(a cherrygo.CreateServer, b cherrygo.CreateServer) bool {
	if a.ProjectID != b.ProjectID || a.Plan != b.Plan || a.Hostname != b.Hostname || a.Image != b.Image ||
		a.Region != b.Region || a.UserData != b.UserData || a.SpotInstance != b.SpotInstance ||
		a.OSPartitionSize != b.OSPartitionSize || a.StorageID != b.StorageID || !slices.Equal(a.SSHKeys, b.SSHKeys) ||
		!slices.Equal(a.IPAddresses, b.IPAddresses) || !maps.Equal(*a.Tags, *b.Tags) {
		return false
	}
	return true

}
func generateCreateRespMock(userDataFile string) func(request *cherrygo.CreateServer) (cherrygo.Server, *cherrygo.Response, error) {
	userData := ""
	if userDataFile != "" {
		userDataRaw, _ := os.ReadFile(userDataFile)
		userData = base64.StdEncoding.EncodeToString(userDataRaw)
	}

	return func(request *cherrygo.CreateServer) (cherrygo.Server, *cherrygo.Response, error) {
		wantReq := cherrygo.CreateServer{
			ProjectID:       1,
			Plan:            "cloud_vps_1",
			Hostname:        "test_hostname",
			Image:           "ubuntu_22_04",
			Region:          "eu_nord_1",
			SSHKeys:         []string{"test_key"},
			IPAddresses:     []string{"test_ip"},
			UserData:        userData,
			Tags:            &map[string]string{"key": "value", "env": "test"},
			SpotInstance:    true,
			OSPartitionSize: 1,
			StorageID:       1,
		}
		if request == nil || !equalCreateServerRequest(wantReq, *request) {
			return cherrygo.Server{}, nil,
				fmt.Errorf("bad creation request, want: %+v, got %+v", &wantReq, request)
		}
		return stubs.Server, nil, nil
	}
}

func setupUserData() string {
	// Set up temporary user data file.
	userDataFile, err := os.CreateTemp("", "test-user-data.*.yaml")
	if err != nil {
		panic(err)
	}
	_, err = userDataFile.Write([]byte("test_user_data"))
	if err != nil {
		panic(err)
	}
	return userDataFile.Name()
}

func destroyUserData(userDataFile string) {
	err := os.Remove(userDataFile)
	if err != nil {
		panic(err)
	}
}

func TestCreateServer(t *testing.T) {
	userDataFile := setupUserData()
	defer destroyUserData(userDataFile)

	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeServersService.CreateResp = generateCreateRespMock(userDataFile)
	})

	f.RootCmd.SetArgs([]string{"server", "create", "--project-id", "1", "--plan", "cloud_vps_1",
		"--hostname", "test_hostname", "--image", "ubuntu_22_04", "--region", "eu_nord_1",
		"--ssh-keys", "test_key", "--ip-addresses", "test_ip", "--userdata-file", userDataFile,
		`--tags=key=value,env=test`, "--spot-instance", "--os-partition-size", "1", "--storage-id", "1"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
	wantTH := []string{"ID", "Name", "Hostname", "Image", "State", "Region"}
	wantTD := [][]string{{"1", "test_name", "test_hostname", "ubuntu_22_04", "deployed", "eu_nord_1"}}
	f.CheckOutput(stubs.Server, wantTH, wantTD)
}

func TestCreateServerAPIError(t *testing.T) {
	userDataFile := setupUserData()
	defer destroyUserData(userDataFile)

	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeServersService.CreateResp = generateCreateRespMock(userDataFile)
	})
	f.RootCmd.SetArgs([]string{"server", "create", "--project-id", "2", "--plan", "cloud_vps_1",
		"--hostname", "test_hostname", "--image", "ubuntu_22_04", "--region", "eu_nord_1",
		"--ssh-keys", "test_key", "--ip-addresses", "test_ip", "--userdata-file", userDataFile,
		`--tags=key=value,env=test`, "--spot-instance", "--os-partition-size", "1", "--storage-id", "1"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not provision a server: bad creation request, want: &{ProjectID:1"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}

func TestCreateServerMissingUserDataError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {})
	f.RootCmd.SetArgs([]string{"server", "create", "--project-id", "1", "--plan", "cloud_vps_1", "--region", "eu_nord_1",
		"--userdata-file", "wrong"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not read userdata-file: open wrong: no such file or directory"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
