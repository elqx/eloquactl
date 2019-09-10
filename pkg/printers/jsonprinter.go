package printers

import (
	"io"
)

type JsonPrinter struct {}

func (p *JsonPrinter) PrintResource(r interface{}, w io.Writer) error {
	// print json implementation
	return nil
}
