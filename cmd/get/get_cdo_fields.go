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
	"context"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/printers"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	exportCdoFieldsLong = templates.LongDesc(`
		Get Eloqua CDO fields and write the result to a file or stdout.

		JSON and CSV file formats are supported`)

	exportCdoFieldsExample = templates.Examples(`
		# Get CDO fields given CDO name.
		eloquactl get cdo-fields NAME

		# Get CDO fields via Bulk API given CDO name.
		eloquactl get cdo-fields NAME --bulk`)
)

type GetCdoFieldsOptions struct {
	Client     func() *bulk.BulkClient
	PrintFlags *cmdutil.PrintFlags

	// Command specific options
	Bulk bool // if True, then feelds are retrieved via Bulk API, not Rest Standard API
}

func NewGetCdoFieldsOptions() *GetCdoFieldsOptions {
	return &GetCdoFieldsOptions{
		Client:     initClient,
		PrintFlags: cmdutil.NewPrintFlags(),
	}
}

func NewCmdGetCdoFields() *cobra.Command {
	o := NewGetCdoFieldsOptions()

	cmd := &cobra.Command{
		Use:     "cdo-fields <NAME>",
		Short:   "Get CDO fields and write the result to a file or stdout.",
		Aliases: []string{"cdo-field"},
		Long:    exportCdoFieldsLong,
		Example: exportCdoFieldsExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.Er("To get CDO fields you need to provide CDO name or ID for the command.")
			}

			//		o.Complete(cmd)
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

func (p *GetCdoFieldsOptions) Validate(cmd *cobra.Command) error {
	// validate shared flags
	if err := p.PrintFlags.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *GetCdoFieldsOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := p.Client()

	var parentId int
	// check if args[0] is numeric or string
	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		// name is given, find id
		// TODO: check cache first)
		cdos, err := client.Cdos.List(ctx)
		if err != nil {
			cmdutil.Er("Failed listing CDOs")
		}

		for _, cdo := range cdos.Items {
			if cdo.Name == args[0] {
				// skipping /customObjects/ prefix
				parentId, err = strconv.Atoi(cdo.Uri[15:])
				if err != nil {
					cmdutil.Er("Error extracting custom object id")
				}
				break
			}
		}
	}

	fields, err := client.Cdos.ListFields(ctx, parentId)
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
