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
	exportAccountsLong = templates.LongDesc(`
		Export Eloqua accounts to a file or stdout.
		
		JSON and CSV file formats are supported`)

	exportAccountsExample = templates.Examples(`
		# Export accounts that were created at and after 2019-01-01
		eloquactl export accounts --filter='{{Account.CreatedAt}}>=2019-01-01'

		# Export specific accounts
		eloquactl export accounts --email-addresses=test1@test.com,test2@test.com'
		
		# Export specific account fields
		eloquactl export accounts --email-addresses=test1@test.com --fields=FirstName:{{Account.Field(C_FirstName)}},LastName:{{Account.Fields(C_LastName)}}`)
)

type ExportAccountsOptions struct {
	Client func() *bulk.BulkClient

	ExportFlags *cmdutil.ExportFlags
	PrintFlags  *cmdutil.PrintFlags

	// Accounts' export specific options
	EmailAddresses []string
	CreatedAt      string
	CreatedAfter   string
	UpdatedAt      string
	UpdatedAfter   string
}

func NewExportAccountsOptions() *ExportAccountsOptions {
	return &ExportAccountsOptions{
		Client:      initClient,
		ExportFlags: cmdutil.NewExportFlags(),
		PrintFlags:  cmdutil.NewPrintFlags(),
	}
}

func NewCmdExportAccounts() *cobra.Command {
	o := NewExportAccountsOptions()
	cmd := &cobra.Command{
		Use:     "accounts",
		Aliases: []string{"account"},
		Short:   "Export Eloqua accounts to a file or stdout",
		Long:    exportAccountsLong,
		Example: exportAccountsExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd)
		},
	}

	// Add shared flags
	o.ExportFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	// Add flags specific to accounts export
	cmd.Flags().StringVar(&o.CreatedAt, "created-at", "", "The date when the account was created.")
	cmd.Flags().StringVar(&o.CreatedAfter, "created-after", "", "The date when the account was created.")
	cmd.Flags().StringVar(&o.UpdatedAt, "updated-at", "", "The date when the account was updated.")
	cmd.Flags().StringVar(&o.UpdatedAfter, "updated-after", "", "The date when the account was updatd.")

	return cmd
}

func (o *ExportAccountsOptions) Complete(cmd *cobra.Command) {
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

func (o *ExportAccountsOptions) Validate() error {
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

func (o *ExportAccountsOptions) Run(cmd *cobra.Command) error {
	ctx := context.Background()
	client := o.Client()
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	var keys []string
	if len(*o.ExportFlags.Fields) == 0 {
		// get fields via api and assign
		// TODO: default fields should be cached
		r, err := client.Accounts.GetFields(ctx)
		if err != nil {
			fmt.Println("Failed to list account fields")
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

	e, err = client.Accounts.CreateExport(ctx, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}
