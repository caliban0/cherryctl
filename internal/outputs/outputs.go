package outputs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

type Outputer interface {
	Output(interface{}, []string, *[][]string) error
	SetFormat(Format)
	Outputln(a ...interface{})
}

type Standard struct {
	Format Format
}

var _ Outputer = (*Standard)(nil)

func outputJSON(in interface{}) error {
	output, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func outputYAML(in interface{}) error {
	output, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (o *Standard) Output(in interface{}, header []string, data *[][]string) error {
	if o.Format == FormatJSON {
		return outputJSON(in)
	} else if o.Format == FormatYAML {
		return outputYAML(in)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeader(header)
		table.AppendBulk(*data)
		table.Render()
		return nil
	}
}

// Outputln outputs lines to standard output by wrapping `fmt.Println(a... any)`.
func (o *Standard) Outputln(a ...interface{}) {
	fmt.Println(a...)
}

func (o *Standard) SetFormat(fmt Format) {
	o.Format = fmt
}
