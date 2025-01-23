package teams_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func TestCreateStorageKeyMissingFlagsError(t *testing.T) {
	for _, v := range [][]string{
		{"team", "create", "--type", "personal"},
		{"team", "create", "--name", "test"},
	} {
		testutils.VerifyMissingFlagsError(t, v)
	}
}

func createRespMock(request *cherrygo.CreateTeam) (cherrygo.Team, *cherrygo.Response, error) {
	wantReq := cherrygo.CreateTeam{
		Name:     "test_team",
		Type:     "personal",
		Currency: "EUR",
	}
	if request == nil || wantReq != *request {
		return cherrygo.Team{}, nil,
			fmt.Errorf("bad creation request, want: %v, got %v", &wantReq, request)
	}
	return stubs.Team, nil, nil
}

func TestCreateTeam(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"team", "create", "--name", "test_team", "--type", "personal", "--currency", "EUR"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	wantTH := []string{"ID", "Name", "Remaining credit", "Hourly usage", "Currency"}
	wantTD := [][]string{{"1", "test_team", "2.200000", "1.100000", "EUR"}}
	f.CheckOutput(stubs.Team, wantTH, wantTD)
}

func TestCreateTeamAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.CreateResp = createRespMock
	})
	f.RootCmd.SetArgs([]string{"team", "create", "--name", "wrong", "--type", "wrong", "--currency", "wrong"})
	err := f.RootCmd.Execute()
	wantMsg := "Could not create a Team: bad creation request," +
		" want: &{test_team personal EUR}, got &{wrong wrong wrong}"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
