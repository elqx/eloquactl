package printers

import (
	"encoding/json"
	"io"
)

type JsonPrinter struct{}

func (p *JsonPrinter) PrintResource(r interface{}, w io.Writer) error {
	data, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}
