package testutils

import (
	"regexp"
	"testing"
)

func VerifyErrorMsg(t *testing.T, got error, want string) {
	r, _ := regexp.Compile(want)
	if got == nil || !r.MatchString(got.Error()) {
		t.Errorf("want err: {%v}, got err: {%v}", want, got)
	}
}

func VerifyMissingFlagsError(t *testing.T, args []string) {
	f := NewExperimentalFixture(t, func(fakeServicer *FakeServicer) {})
	f.RootCmd.SetArgs(args)
	err := f.RootCmd.Execute()
	const wantMsg = `required flag\(s\) "(.*)" not set`
	VerifyErrorMsg(t, err, wantMsg)
}
