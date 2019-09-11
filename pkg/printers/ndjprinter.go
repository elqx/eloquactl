package printers

import (
	"fmt"
	"io"
	"encoding/json"

	"github.com/elqx/eloqua-go/eloqua/bulk"
)

type NdjPrinter struct {}

func (p *NdjPrinter) PrintResource(r interface{}, w io.Writer) error {
	items, ok := r.([]bulk.Item)

	if ok {
		for _, item := range items {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(w, string(b))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
