package outputs

import (
	"encoding/json"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/term"
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
}

type Standard struct {
	Format Format
}

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

// Not sure how to handle cases where the table doesn't fit into the terminal, of if there
// even is a way to do that cleanly.
//
// Using just a global limit with cfg.MaxWidth doesn't
// work, because it divides the available space to each column equally, so long columns
// get breaks, even if there was enough space on the screen to render the whole table.
//
// So that leaves us with a per-column width setting. If the per-column width exceeds the
// available screen space, it gets wrapped into the next line, along with the table borders.
// Using cfg.MaxWidth in conjunction with per-column width doesn't work, so it would
// be necessary to keep track of the available screen "budget",
// taking into account borders ant padding.
//
// In which case, we could either: a) shrink columns until everything fits, or
// b) drop the last columns, with the assumption that the first ones are more important.
// I think I would prefer option b, since if we WrapBreak things like UUIDs, it might
// be hard to see the little ↩ symbol. It could also break scripts relying on `cherryctl`.
//
// We could also keep everything as is :).


func colWidths(header []string, data [][]string) (m tw.Mapper[int, int], total int) {
	widths := tw.NewMapper[int, int]()
	for i, v := range header {
		c := utf8.RuneCountInString(v) + 2
		widths.Set(i, c)
		total += c
	}
	for i, row := range data {
		if len(row) != len(header) {
			panic(fmt.Sprintf(
				"header has %d columns, row [%d] has %d columns",
				len(header), i, len(row)))
		}
		for j, v := range row {
			if c := utf8.RuneCountInString(v) + 2; c > widths.Get(j) {
				total -= widths.Get(j)
				widths.Set(j, c)
				total += c
			}
		}
	}
	return widths, total
}

func outputTable(header []string, data [][]string) error {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 120
	}

	_, total := colWidths(header, data)

	table := tablewriter.NewWriter(os.Stdout)
	table = table.Configure(func(cfg *tablewriter.Config) {
		cfg.Row.Alignment.Global = tw.AlignLeft

		if total > termWidth - 2 {
			cfg.MaxWidth = termWidth - 2
		}
		//cfg.MaxWidth = termWidth - 2
		//cfg.Widths.PerColumn = colWidths(header, data)
		cfg.Row.Formatting.AutoWrap = tw.WrapBreak
	})

	table.Header(header)
	if err := table.Bulk(data); err != nil {
		return err
	}
	return table.Render()
}

func (o *Standard) Output(in interface{}, header []string, data *[][]string) error {
	if o.Format == FormatJSON {
		return outputJSON(in)
	} else if o.Format == FormatYAML {
		return outputYAML(in)
	} else {
		return outputTable(header, *data)
	}
}

func (o *Standard) SetFormat(fmt Format) {
	o.Format = fmt
}
