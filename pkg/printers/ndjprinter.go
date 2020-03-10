package printers

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloqua-go/eloqua/rest"
)

type NdjPrinter struct{}

func (p *NdjPrinter) PrintResource(r interface{}, w io.Writer) error {
	//iitems, ok := r.([]bulk.Item)

	switch r := r.(type) {
	case []bulk.Item:
		for _, item := range r {
			printItem(&w, item)
		}
	case []bulk.CdoField:
		for _, item := range r {
			printItem(&w, item)
		}
	case []bulk.ContactField:
		for _, item := range r {
			printItem(&w, item)
		}
	case []rest.Campaign:
		for _, item := range r {
			printItem(&w, item)
		}
	case []rest.Form:
		for _, item := range r {
			printItem(&w, item)
		}
	case []rest.Email:
		for _, item := range r {
			printItem(&w, item)
		}
	}

	return nil
}

func printItem(w *io.Writer, item interface{}) error {
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(*w, string(b))
	if err != nil {
		return err
	}
	return nil
}
