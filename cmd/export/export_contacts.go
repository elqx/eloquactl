package export

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/elqx/eloqua-go/eloqua/bulk"
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
	PrintFlags *cmdutil.PrintFlags
	NoHeaders bool
	OutputFormat string
	Sort bool

	UTC bool
	AutoDeleteDuration string
	DataRetentionDuration string
	Name string
	Fields string
	Filter string // raw filter should not be exposed to the user, but it can be in this struct (flags are exposed to the user, not this struct)
	MaxRecords uint
	// contact filters are the fields
	EmailAddresses []string
	CreatedAt string
	CreatedAfter string
	UpdatedAt string
	UpdatedAfter string
}

func NewExportContactsOptions() *ExportContactsOptions {
	return &ExportContactsOptions{
		Client: initClient,
		PrintFlags: cmdutil.NewPrintFlags(),
	}
}

func NewCmdExportContacts() *cobra.Command {
	o := NewExportContactsOptions()

	cmd := &cobra.Command{
		Use: "contacts",
		Aliases: []string{"contact"},
		Short: "Export Eloqua contacts to a file or stdout",
		Long: exportContactsLong,
		Example: exportContactsExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd)
		},
	}
	cmdutil.AddStagingFlags(cmd)
	cmdutil.AddNameFlag(cmd)
	cmd.Flags().BoolP("utc", "u", true, "Whether or not system timestamps will be exported in UTC.")
	cmd.Flags().String("fields", "", "List of fields to be included in the export operation.")
	cmd.Flags().Int("max-records", 1000, "The maximum amount of records.")
	cmd.Flags().String("created-at", "", "The date when the contact was created.")
	cmd.Flags().String("created-after", "", "The date when the contact was created.")
	cmd.Flags().String("updated-at", "", "The date when the contact was updated.")
	cmd.Flags().String("updated-after", "", "The date when the contact was updatd.")
	cmd.Flags().StringSlice("email-addresses", []string{}, "Contacts' email addresses.")
	//cmd.Flags().StringP("filter", "f", "", "Contact filter. EML statement")
	return cmd
}

func (p *ExportContactsOptions) Complete(cmd *cobra.Command) error {
	p.AutoDeleteDuration = cmdutil.GetFlagString(cmd, "auto-delete-duration")
	p.DataRetentionDuration = cmdutil.GetFlagString(cmd, "data-retention-duration")

	p.UTC = cmdutil.GetFlagBool(cmd, "utc")

	// name does not have a default, generate name if not specified
	p.Name = cmdutil.GetFlagString(cmd, "name")
	if len(p.Name) == 0 {
		p.Name = generateName()
	}

	p.Fields = cmdutil.GetFlagString(cmd, "fields")
	p.MaxRecords = cmdutil.GetFlagUint(cmd, "max-records")
	p.CreatedAt = cmdutil.GetFlagString(cmd, "created-at")
	p.CreatedAfter = cmdutil.GetFlagString(cmd, "created-after")
	p.UpdatedAt = cmdutil.GetFlagString(cmd, "updated-at")
	p.UpdatedAfter = cmdutil.GetFlagString(cmd, "updated-after")
	p.EmailAddresses = cmdutil.GetFlagStringSlice(cmd, "email-addresses")

	if len(p.EmailAddresses) > 0 {
		// add email addresses to filter
	}

	if len(p.CreatedAt) > 0{
		// add CreatedAt to the filter
	}

	if len(p.CreatedAfter) > 0 {
		// add CreatedAt > CreatedAfter to the filter
	}

	if len(p.UpdatedAt) > 0 {
		// add UpdatedAt to the filter
	}

	if len(p.UpdatedAfter) > 0 {
		// add UpdatedAt > UpdatedAfter to the filter
	}

	return nil
}

func (p *ExportContactsOptions) Validate() error {
	// check that staging options follow ISO-8601 standard
	if err := checkISO8601(p.AutoDeleteDuration); err != nil {
		fmt.Println("Failed validating auto-delete-duration. Value should follow ISO-8601 standard.")
	}

	if err := checkISO8601(p.DataRetentionDuration); err != nil {
		fmt.Println("Failed validation data-retention-duration. Value should follow ISO-8601 standard.")
	}

	// check that name is no longer than 100 (based on doc)
	if len(p.Name) > 100 {
		fmt.Println("Export name must be no longer than 100 characters long.")
	}

	// validate fields if they were provided by the user
	if len(p.Fields) > 0 {
		// TODO: check regex expression
	}

	if p.MaxRecords < 0 {
		fmt.Println("max-records should be greater than zero")
	}

	if err := checkDate(p.CreatedAt); err != nil {
		return err
	}

	if err := checkDate(p.CreatedAfter); err != nil {
		return err
	}

	if err := checkDate(p.UpdatedAt); err != nil {
		return err
	}

	if err := checkDate(p.UpdatedAfter); err != nil {
		return err
	}

	// TODO: check email addresses using regex
	return nil
}

func (p *ExportContactsOptions) Run(cmd *cobra.Command) error {
	ctx := context.Background()
	client := p.Client()
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	var keys []string
	if len(p.Fields) == 0 {
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
		k, err := parseFieldsStr(p.Fields, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys = k
	}

	if len(p.Filter) == 0 {
		// get fields via api and construct the filter
		// fields should be cached
	}

	printer, err := p.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	// should have Filter struct in the client library
	//var filter strings.Builder
	e := &bulk.Export{
		AreSystemTimestampsInUTC: p.UTC,
		AutoDeleteDuration: p.AutoDeleteDuration,
		DataRetentionDuration: p.DataRetentionDuration,
		Name: p.Name,
		Fields: fields,
		Filter: "'{{Contact.CreatedAt}}' >= '2019-09-05'",
		//Filter: p.Filter,
		MaxRecords: p.MaxRecords,
	}

	e, err = client.Contacts.CreateExport(ctx, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}

