package printers

import (
	"io"
)

type ResourcePrinter interface {
	PrintResource(interface{}, io.Writer) error
}

type ResourcePrinterFunc func(interface{}, io.Writer) error

func (fn ResourcePrinterFunc) PrintResource(r interface{}, w io.Writer) error {
	return fn(r, w)
}
