package teams_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func getResp(teamID int, opts *cherrygo.GetOptions) (cherrygo.Team, *cherrygo.Response, error) {
	if teamID != 1 {
		return cherrygo.Team{}, nil, fmt.Errorf("team id does not match: want %d, got %d", 1, teamID)
	}
	return stubs.Team, nil, nil
}

func TestGetTeam(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.GetResp = getResp
	})
	f.RootCmd.SetArgs([]string{"team", "get", "1"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	wantTH := []string{"ID", "Name", "Remaining credit", "Hourly usage", "Currency"}
	wantTD := [][]string{{"1", "test_team", "2.200000", "1.100000", "EUR"}}
	f.CheckOutput(stubs.Team, wantTH, wantTD)
}

func TestGetTeamAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.GetResp = getResp
	})
	f.RootCmd.SetArgs([]string{"team", "get", "2"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get team: team id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
