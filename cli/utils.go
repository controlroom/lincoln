package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var yellowOut = color.New(color.FgYellow).SprintFunc()

type Table struct {
	Table *tablewriter.Table
}

func getTable(headers ...string) *Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetBorder(false)

	table.SetHeader(headers)

	return &Table{table}
}

func (t *Table) AppendLine(line ...string) {
	t.Table.Append(line)
}

func (t *Table) Render() {
	t.Table.Render()
}
