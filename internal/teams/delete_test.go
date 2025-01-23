package teams_test

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/testutils"
	"github.com/cherryservers/cherrygo/v3"
	"testing"
)

func deleteResp(teamID int) (*cherrygo.Response, error) {
	if teamID != 1 {
		return nil, fmt.Errorf("team id does not match: want %d, got %d", 1, teamID)
	}
	return nil, nil
}

func TestDeleteTeamWithForce(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.DeleteResp = deleteResp
	})
	f.RootCmd.SetArgs([]string{"team", "delete", "1", "-f"})
	if err := f.RootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}

	f.CheckOutputLine(fmt.Sprintln("Team", 1, "successfully deleted."))
}

func TestDeleteTeamWithForceAPIError(t *testing.T) {
	f := testutils.NewExperimentalFixture(t, func(fakeServicer *testutils.FakeServicer) {
		fakeServicer.FakeTeamsService.DeleteResp = deleteResp
	})
	f.RootCmd.SetArgs([]string{"team", "delete", "2", "--force"})
	err := f.RootCmd.Execute()
	const wantMsg = "Could not delete a Team: team id does not match: want 1, got 2"
	testutils.VerifyErrorMsg(t, err, wantMsg)
}
