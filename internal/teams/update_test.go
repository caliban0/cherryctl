package teams_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"reflect"
	"testing"
)

func updateRespMock(teamID int, request *cherrygo.UpdateTeam) (cherrygo.Team, *cherrygo.Response, error) {
	// Request fields are, unfortunately, pointers.
	wName, wType, wCurrency := "test_team", "personal", "EUR"
	wantReq := cherrygo.UpdateTeam{
		Name:     &wName,
		Type:     &wType,
		Currency: &wCurrency,
	}
	// Need to use deep equal because of the request pointer types.
	if request == nil || !reflect.DeepEqual(request, &wantReq) {
		return cherrygo.Team{}, nil,
			fmt.Errorf("bad update request, want: %v, got %v", &wantReq, request)
	}
	return stubs.Team, nil, nil
}

func TestUpdateTeam(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.UpdateResp = updateRespMock
	})
	f.RootCmd.SetArgs([]string{"team", "update", "--name", "test_team", "--type", "personal", "--currency", "EUR"})

	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	wantTH := []string{"ID", "Name", "Remaining credit", "Hourly usage", "Currency"}
	wantTD := [][]string{{"1", "test_team", "2.200000", "1.100000", "EUR"}}
	f.CheckOutput(stubs.Team, wantTH, wantTD)
}

func TestUpdateTeamAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.UpdateResp = updateRespMock
	})
	f.RootCmd.SetArgs([]string{"team", "update", "--name", "wrong", "--type", "wrong", "--currency", "wrong"})
	err := f.RootCmd.Execute()

	// Check only prefix, because of pointer types in request.
	wantMsg := "Could not update team: bad update request, want: "
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
