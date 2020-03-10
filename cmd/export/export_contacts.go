package export

import (
	"context"
	"fmt"
	"os"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
)

var (
	exportContactsLong = templates.LongDesc(`
		Export Eloqua contacts to a file or stdout.
		
		JSON and CSV file formats are supported`)

	exportContactsExample = templates.Examples(`
		# Export contacts that were created at and after 2019-01-01
		eloquactl export contacts --filter='{{Contact.CreatedAt}}>=2019-01-01'

		# Export specific contacts
		eloquactl export contacts --email-addresses=test1@test.com,test2@test.com'
		
		# Export specific contact fields
		eloquactl export contacts --email-addresses=test1@test.com --fields=FirstName:{{Contact.Field(C_FirstName)}},LastName:{{Contact.Fields(C_LastName)}}`)
)

type ExportContactsOptions struct {
	Client func() *bulk.BulkClient

	ExportFlags *cmdutil.ExportFlags
	PrintFlags  *cmdutil.PrintFlags

	// Contacts' export specific options
	EmailAddresses []string
	CreatedAt      string
	CreatedAfter   string
	UpdatedAt      string
	UpdatedAfter   string
}

func NewExportContactsOptions() *ExportContactsOptions {
	return &ExportContactsOptions{
		Client:      initClient,
		ExportFlags: cmdutil.NewExportFlags(),
		PrintFlags:  cmdutil.NewPrintFlags(),
	}
}

func NewCmdExportContacts() *cobra.Command {
	o := NewExportContactsOptions()
	cmd := &cobra.Command{
		Use:     "contacts",
		Aliases: []string{"contact"},
		Short:   "Export Eloqua contacts to a file or stdout",
		Long:    exportContactsLong,
		Example: exportContactsExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd)
		},
	}

	// Add shared flags
	o.ExportFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	// Add flags specific to contacts export
	cmd.Flags().StringVar(&o.CreatedAt, "created-at", "", "The date when the contact was created.")
	cmd.Flags().StringVar(&o.CreatedAfter, "created-after", "", "The date when the contact was created.")
	cmd.Flags().StringVar(&o.UpdatedAt, "updated-at", "", "The date when the contact was updated.")
	cmd.Flags().StringVar(&o.UpdatedAfter, "updated-after", "", "The date when the contact was updatd.")
	cmd.Flags().StringSliceVar(&o.EmailAddresses, "email-addresses", []string{}, "Contacts' email addresses.")

	return cmd
}

func (o *ExportContactsOptions) Complete(cmd *cobra.Command) {
	// StagingFlags, ExportFlags and PrintFlags are completed
	// here should only be the completion of the filter option, (and name)
	if len(o.EmailAddresses) > 0 {
		// add email addresses to filter
	}

	if len(o.CreatedAt) > 0 {
		// add CreatedAt to the filter
	}

	if len(o.CreatedAfter) > 0 {
		// add CreatedAt > CreatedAfter to the filter
	}

	if len(o.UpdatedAt) > 0 {
		// add UpdatedAt to the filter
	}

	if len(o.UpdatedAfter) > 0 {
		// add UpdatedAt > UpdatedAfter to the filter
	}
}

func (o *ExportContactsOptions) Validate() error {
	//var e Errors
	if err := o.ExportFlags.Validate(); err != nil {
		//aggregateError(err)
		return err
	}

	if err := o.PrintFlags.Validate(); err != nil {
		//aggregateError(err)
		return err
	}

	if err := checkDate(o.CreatedAt); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := checkDate(o.CreatedAfter); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := checkDate(o.UpdatedAt); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := checkDate(o.UpdatedAfter); err != nil {
		//e = aggregateError(err)
		return err
	}

	return nil
}

func (o *ExportContactsOptions) Run(cmd *cobra.Command) error {
	ctx := context.Background()
	client := o.Client()
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	var keys []string
	if len(*o.ExportFlags.Fields) == 0 {
		// get fields via api and assign
		// TODO: default fields should be cached
		r, err := client.Contacts.GetFields(ctx)
		if err != nil {
			fmt.Println("Failed to list contact fields")
			os.Exit(1)
		}

		for _, f := range r.Items {
			fields[f.InternalName] = f.Statement
		}
	} else {
		k, err := parseFieldsStr(*o.ExportFlags.Fields, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys = k
	}

	if len(*o.ExportFlags.Filter) == 0 {
		// get fields via api and construct the filter
		// fields should be cached
	}

	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	// should have Filter struct in the client library
	//var filter strings.Builder
	e := &bulk.Export{
		AreSystemTimestampsInUTC: *o.ExportFlags.AreSystemTimestampsInUTC,
		AutoDeleteDuration:       *o.ExportFlags.StagingFlags.AutoDeleteDuration,
		DataRetentionDuration:    *o.ExportFlags.StagingFlags.DataRetentionDuration,
		Name:                     *o.ExportFlags.Name,
		Fields:                   fields,
		Filter:                   *o.ExportFlags.Filter,
		//Filter: p.Filter,
		MaxRecords: *o.ExportFlags.MaxRecords,
	}

	e, err = client.Contacts.CreateExport(ctx, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}
