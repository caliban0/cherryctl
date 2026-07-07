package plans

import (
	"fmt"
	"testing"

	"github.com/cherryservers/cherryctl/internal/fakes"
	"github.com/cherryservers/cherrygo/v3"
)

func TestGet(t *testing.T) {
	cases := []struct {
		title            string
		args             []string
		getOpts          *cherrygo.GetOptions
		wantClientParams []any
		wantTd           [][]string
	}{
		{
			title:            "only plan slug",
			args:             []string{"test-plan"},
			wantClientParams: []any{"test-plan", (*cherrygo.GetOptions)(nil)},
			wantTd:           [][]string{{"test-plan", "test-region", "1", fmt.Sprintf("%f", 1.0), "2", fmt.Sprintf("%f", 0.5)},
				{"test-plan", "test-region-2", "1", fmt.Sprintf("%f", 1.0), "2", fmt.Sprintf("%f", 0.5)}},
		},
		{
			title:            "filter by region",
			args:             []string{"test-plan", "--region", "test-region"},
			wantClientParams: []any{"test-plan", (*cherrygo.GetOptions)(nil)},
			wantTd: [][]string{{"test-plan", "test-region", "1", fmt.Sprintf("%f", 1.0), "2", fmt.Sprintf("%f", 0.5)}},
		},
		{
			title:            "shorthands",
			args:             []string{"test-plan", "-r", "test-region"},
			wantClientParams: []any{"test-plan", (*cherrygo.GetOptions)(nil)},
			wantTd: [][]string{{"test-plan", "test-region", "1", fmt.Sprintf("%f", 1.0), "2", fmt.Sprintf("%f", 0.5)}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			fakeSvc := fakes.PlanService{}
			fakeOut := fakes.Outputer{}
			dep := fakeDeps{
				svc:  &fakeSvc,
				out:  &fakeOut,
				opts: tc.getOpts,
			}

			c := Command{
				Deps: dep,
			}
			cmd := c.get()
			cmd.SetArgs(tc.args)
			cmd.SilenceUsage = true

			err := cmd.Execute()
			if err != nil {
				t.Fatal(err.Error())
			}

			if len(fakeSvc.Calls) != 1 {
				t.Fatalf("want 1 api call, got %d", len(fakeSvc.Calls))
			}
			fakeSvc.Calls[0].AssertMethod(t, "GetBySlug")
			fakeSvc.Calls[0].AssertParams(t, tc.wantClientParams...)

			if len(fakeOut.Calls) != 1 {
				t.Fatalf("want 1 output call, got %d", len(fakeOut.Calls))
			}

			wantPlan := fakeSvc.Plan()
			wantTh := []string{"plan slug", "region slug", "stock hourly", "hourly price", "stock spot", "spot price"}
			fakeOut.Calls[0].Assert(t, wantPlan, wantTh, tc.wantTd)
		})
	}
}
