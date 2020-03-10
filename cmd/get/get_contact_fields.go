/*
Copyright © 2019 elqx <ignotas.petrulis@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package get

import (
	"os"
	//"fmt"
	"context"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/printers"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
)

var (
	getContactFieldsLong = templates.LongDesc(`
		Get Eloqua Contact fields and write the result to a file or stdout.

		JSON and CSV file formats are supported`)

	getContactFieldsExample = templates.Examples(`
		# Get Contact fields
		eloquactl get contact-fields`)
)

type GetContactFieldsOptions struct {
	Client     func() *bulk.BulkClient
	PrintFlags *cmdutil.PrintFlags

	// Command specific options
	Bulk bool // if True, then feelds are retrieved via Bulk API, not Rest Standard API
}

func NewGetContactFieldsOptions() *GetContactFieldsOptions {
	return &GetContactFieldsOptions{
		Client:     initClient,
		PrintFlags: cmdutil.NewPrintFlags(),
	}
}

func NewCmdGetContactFields() *cobra.Command {
	o := NewGetContactFieldsOptions()

	cmd := &cobra.Command{
		Use:     "contact-fields",
		Short:   "Get Contact fields and write the result to a file or stdout.",
		Aliases: []string{"cdo-field"},
		Long:    getContactFieldsLong,
		Example: getContactFieldsExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate(cmd)
			o.Run(cmd, args)
		},
	}

	// Add shared flags
	o.PrintFlags.AddFlags(cmd)
	// Add command specific flags
	cmd.Flags().BoolVar(&o.Bulk, "bulk", false, "Specifies whether result should be retrieved via Bulk API.")

	return cmd
}

func (p *GetContactFieldsOptions) Validate(cmd *cobra.Command) error {
	// validate shared flags
	if err := p.PrintFlags.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *GetContactFieldsOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := p.Client()

	fields, err := client.Contacts.GetFields(ctx)
	if err != nil {
		return err
	}

	printer, err := p.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	w := printers.NewTabWriter(os.Stdout)
	printer.PrintResource(fields.Items, w)
	w.Flush()

	return nil
}
