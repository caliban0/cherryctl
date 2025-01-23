package testutils

import (
	"fmt"
	"github.com/cherryservers/cherryctl/internal/outputs"
)

var _ outputs.Outputer = &SpyOutputer{}

// SpyOutputer is an `outputer` implementation that monitors and checks
// output data for testing purposes.
type SpyOutputer struct {
	Format      outputs.Format
	Resource    interface{}
	TableHeader []string
	TableData   [][]string
	Line        string
}

// Output stores output data for later testing.
func (s *SpyOutputer) Output(in interface{}, header []string, data *[][]string) error {
	s.Resource = in
	s.TableHeader = header
	s.TableData = *data
	return nil
}

// SetFormat sets output format. This method is required by the `outputer` interface,
// but is not used when testing.
func (s *SpyOutputer) SetFormat(format outputs.Format) {
	s.Format = format
}

// Outputln stores the output line as a string for later testing.
// Most commands don't output lines, the exception being commands
// that don't have a response body, such as `delete`.
func (s *SpyOutputer) Outputln(a ...interface{}) {
	s.Line = fmt.Sprintln(a...)
}
