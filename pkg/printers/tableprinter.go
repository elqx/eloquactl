package printers

import (
	"io"
	"fmt"
	"strings"

	"github.com/elqx/eloqua-go/eloqua/bulk"
)

type TablePrinter struct {}

func (p *TablePrinter) PrintResource(r interface{}, w io.Writer) error {
	table, ok := r.([]bulk.Item)
	var headers []string
	var keys []string

	if ok {
		row := table[0]

		for cellName, _ := range row {
			keys = append(keys, cellName)
			headers = append(headers, strings.ToUpper(cellName))
		}

		printHeader(headers, w)

		for _, row := range  table {
			for _, key := range keys {
				if value, exists := row[key]; exists {
					fmt.Fprintf(w, "%s\t", value)
				}
			}
			fmt.Fprint(w, "\n")
		}
	}
	return nil
}


func printHeader(columnNames []string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t")); err != nil {
		return err
	}
	return nil
}
