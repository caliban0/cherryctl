package teams_test

import (
	"errors"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherryctl/internal/testutils/stubs"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func listResp(opts *cherrygo.GetOptions) ([]cherrygo.Team, *cherrygo.Response, error) {
	return []cherrygo.Team{stubs.Team}, nil, nil
}

func listRespErr(opts *cherrygo.GetOptions) ([]cherrygo.Team, *cherrygo.Response, error) {
	return []cherrygo.Team{}, nil, errors.New("api_error")
}

func TestListTeams(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.ListResp = listResp
	})
	f.RootCmd.SetArgs([]string{"team", "list"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	wantTH := []string{"ID", "Name", "Remaining credit", "Hourly usage", "Currency"}
	wantTD := [][]string{{"1", "test_team", "2.200000", "1.100000", "EUR"}}
	f.CheckOutput([]cherrygo.Team{stubs.Team}, wantTH, wantTD)
}

func TestListTeamsAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.ListResp = listRespErr
	})
	f.RootCmd.SetArgs([]string{"team", "list"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not get teams list: api_error"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
