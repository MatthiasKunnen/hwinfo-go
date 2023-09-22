package text

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Column struct {
	width int
}

// TablePrinter is a utility to print a text table where the columns are sized automatically to the
// content.
type TablePrinter struct {
	Columns []Column
	Padding string
	data    [][]string
	output  io.Writer
}

func NewTablePrinter(output io.Writer, columns []Column, padding string) *TablePrinter {
	return &TablePrinter{
		Columns: columns,
		Padding: padding,
		data:    make([][]string, 0),
		output:  output,
	}
}

func (receiver *TablePrinter) Append(row []string) {
	receiver.data = append(receiver.data, row)
	for i, cell := range row {
		cellLength := utf8.RuneCountInString(cell)

		if receiver.Columns[i].width < cellLength {
			receiver.Columns[i].width = cellLength
		}
	}
}

func (receiver *TablePrinter) Write() error {
	for _, row := range receiver.data {
		for i, cell := range row {
			if i < len(row)-1 {
				_, err := fmt.Fprintf(
					receiver.output,
					"%-*s%s",
					receiver.Columns[i].width,
					cell,
					receiver.Padding,
				)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintf(receiver.output, "%s", cell)
				if err != nil {
					return err
				}
			}
		}
		_, err := receiver.output.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
