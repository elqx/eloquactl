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
	"github.com/elqx/eloquactl/pkg/printers"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
	"os"

	"github.com/elqx/eloqua-go/eloqua/rest"
)

var (
	getEmailsLong = templates.LongDesc(`
		Get Eloqua emails and write the result to a file or stdout.

		JSON and CSV file formats are supported`)

	getEmailsExample = templates.Examples(`
		# Get email given its name.
		eloquactl get emails NAME`)
)

type GetEmailsOptions struct {
	Client    func() *rest.RestClient
	ListFlags *cmdutil.ListFlags

	PrintFlags *cmdutil.PrintFlags

	All bool
}

func NewGetEmailsOptions() *GetEmailsOptions {
	return &GetEmailsOptions{
		Client:     initRestClient,
		ListFlags:  cmdutil.NewListFlags(),
		PrintFlags: cmdutil.NewPrintFlags(),
	}
}

func NewCmdGetEmails() *cobra.Command {
	o := NewGetEmailsOptions()

	cmd := &cobra.Command{
		Use:     "emails <NAME>",
		Short:   "Get Eloqua emails and write the result to a file or stdout.",
		Aliases: []string{"email"},
		Long:    getEmailsLong,
		Example: getEmailsExample,
		Run: func(cmd *cobra.Command, args []string) {
			//if len(args) < 1 {
			//	cmdutil.Er("To get emails you need to provide its name or ID for the command.")
			//}

			//		o.Complete(cmd)
			o.Validate(cmd)
			o.Run(cmd, args)
		},
	}

	o.ListFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.All, "all", false, "Specifies whether all emails should be retrieved.")

	return cmd
}

func (p *GetEmailsOptions) Validate(cmd *cobra.Command) error {
	// validate shared flags
	if err := p.PrintFlags.Validate(); err != nil {
		return err
	}

	// validate ListFlags

	return nil
}

func (p *GetEmailsOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := p.Client()

	printer, err := p.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	w := printers.NewTabWriter(os.Stdout)

	opts := &rest.GetOptions{
		Count:   *p.ListFlags.Count,
		Depth:   *p.ListFlags.Depth,
		OrderBy: *p.ListFlags.OrderBy,
		Page:    *p.ListFlags.Page,
		Search:  *p.ListFlags.Search,
	}

	if p.All {
		pageSize := *p.ListFlags.Count
		totalResults := 99999999 // upper estimate for the number of emails
		for page := 1; (page-1)*pageSize <= totalResults; page++ {
			opts.Page = page
			emails, err := client.Emails.ListEmails(ctx, opts)
			if err != nil {
				return err
			}
			totalResults = emails.Total
			printer.PrintResource(emails.Elements, w)
			w.Flush()
		}
	} else {
		emails, err := client.Emails.ListEmails(ctx, opts)
		if err != nil {
			return err
		}
		printer.PrintResource(emails.Elements, w)
		w.Flush()
	}

	return nil
}
