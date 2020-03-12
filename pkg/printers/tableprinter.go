package printers

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloqua-go/eloqua/rest"
)

type TablePrinter struct{}

func (p *TablePrinter) PrintResource(r interface{}, w io.Writer) error {
	var headers []string
	var keys []string

	switch r := r.(type) {
	case []bulk.Item:
		row := r[0]
		for cellName, _ := range row {
			keys = append(keys, cellName)
			headers = append(headers, strings.ToUpper(cellName))
		}

		printHeader(headers, w)

		for _, row := range r {
			for _, key := range keys {
				if value, exists := row[key]; exists {
					fmt.Fprintf(w, "%s\t", value)
				}
			}
			fmt.Fprint(w, "\n")
		}
	case []bulk.CdoField:
		for i, row := range r {
			v := reflect.ValueOf(row)

			if i == 0 {
				for j := 0; j < v.NumField(); j++ {
					headers = append(headers, strings.ToUpper(v.Type().Field(j).Name))
				}
				printHeader(headers, w)
			}

			for j := 0; j < v.NumField(); j++ {
				fmt.Fprintf(w, "%v\t", v.Field(j).Interface())
			}
			fmt.Fprint(w, "\n")
		}
	case []bulk.ContactField:
		for i, row := range r {
			v := reflect.ValueOf(row)

			if i == 0 {
				for j := 0; j < v.NumField(); j++ {
					headers = append(headers, strings.ToUpper(v.Type().Field(j).Name))
				}
				printHeader(headers, w)
			}

			for j := 0; j < v.NumField(); j++ {
				fmt.Fprintf(w, "%v\t", v.Field(j).Interface())
			}
			fmt.Fprint(w, "\n")
		}
	case []rest.Campaign:
		fields := []string{"id", "name", "currentStatus", "createdAt", "createdBy", "updatedAt", "updatedBy"}
		for i, c := range r {
			//v := reflect.ValueOf(row)
			if i == 0 {
				for _, field := range fields {
					headers = append(headers, strings.ToUpper(field))
				}
				printHeader(headers, w)
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Id, c.Name, c.CurrentStatus, c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy)
		}
	case []rest.Email:
		fields := []string{"id", "name", "currentStatus", "createdAt", "createdBy", "updatedAt", "updatedBy"}
		for i, c := range r {
			//v := reflect.ValueOf(row)
			if i == 0 {
				for _, field := range fields {
					headers = append(headers, strings.ToUpper(field))
				}
				printHeader(headers, w)
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Id, c.Name, c.CurrentStatus, c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy)
		}
	case []rest.Form:
		fields := []string{"id", "name", "currentStatus", "createdAt", "createdBy", "updatedAt", "updatedBy"}
		for i, c := range r {
			//v := reflect.ValueOf(row)
			if i == 0 {
				for _, field := range fields {
					headers = append(headers, strings.ToUpper(field))
				}
				printHeader(headers, w)
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Id, c.Name, c.CurrentStatus, c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy)
		}
	case []rest.EmailGroup:
		fields := []string{"id", "name", "currentStatus", "createdAt", "createdBy", "updatedAt", "updatedBy"}
		for i, c := range r {
			if i == 0 {
				for _, field := range fields {
					headers = append(headers, strings.ToUpper(field))
				}
				printHeader(headers, w)
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n", c.Id, c.Name, c.CurrentStatus, c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy)
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
