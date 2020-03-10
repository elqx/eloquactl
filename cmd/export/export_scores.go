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
	exportScoresLong = templates.LongDesc(`
		Export Eloqua lead scores to a file or stdout.
		
		JSON and CSV file formats are supported`)

	exportScoresExample = templates.Examples(`
		# Export lead scores of a lead scoring model named 'model1'
		eloquactl export scores model1 --output=ndj
		
		# Export lead scores of a lead scoring models with id 1.
		eloquactl export scores 1 --output=ndj`)
)

type ExportScoresOptions struct {
	Client func() *bulk.BulkClient

	ExportFlags *cmdutil.ExportFlags
	PrintFlags  *cmdutil.PrintFlags
}

func NewExportScoresOptions() *ExportScoresOptions {
	return &ExportScoresOptions{
		Client:      initClient,
		ExportFlags: cmdutil.NewExportFlags(),
		PrintFlags:  cmdutil.NewPrintFlags(),
	}
}

func NewCmdExportScores() *cobra.Command {
	o := NewExportScoresOptions()

	cmd := &cobra.Command{
		Use:     "scores",
		Aliases: []string{"score"},
		Short:   "Export Eloqua lead scores to a file or stdout",
		Long:    exportScoresLong,
		Example: exportScoresExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd, args)
		},
	}

	// Add shared flags
	o.ExportFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	return cmd
}

func (o *ExportScoresOptions) Complete(cmd *cobra.Command) {

}

func (o *ExportScoresOptions) Validate() error {
	//var e Errors
	if err := o.ExportFlags.Validate(); err != nil {
		//aggregateError(err)
		return err
	}

	if err := o.PrintFlags.Validate(); err != nil {
		//aggregateError(err)
		return err
	}

	return nil
}

func (o *ExportScoresOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := o.Client()
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	var keys []string
	if len(*o.ExportFlags.Fields) == 0 {
		// get fields via api and assign
		// TODO: default fields should be cached
		r, err := client.Contacts.ListLeadModels(ctx)
		if err != nil {
			fmt.Println("Failed to list lead scoring models")
			os.Exit(1)
		}

		for _, model := range r.Items {
			if model.Name == args[0] {
				for _, f := range model.Fields {
					fields[f.Name] = f.Statement
				}
				break
			}
		}

		if len(fields) == 0 {
			fmt.Printf("Lead Scoring Model %v does not exist", args[0])
			os.Exit(1)
		}
	} else {
		k, err := parseFieldsStr(*o.ExportFlags.Fields, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys = k
	}

	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	// add {{Contact.Id}} and {{Contact.Field(C_EmailAddress)}} fields to the export
	fields["ContactId"] = "{{Contact.Id}}"
	fields["EmailAddress"] = "{{Contacts.Fields(C_EmailAddress)}}"
	// should have Filter struct in the client library
	//var filter strings.Builder
	e := &bulk.Export{
		AreSystemTimestampsInUTC: *o.ExportFlags.AreSystemTimestampsInUTC,
		AutoDeleteDuration:       *o.ExportFlags.StagingFlags.AutoDeleteDuration,
		DataRetentionDuration:    *o.ExportFlags.StagingFlags.DataRetentionDuration,
		Name:                     *o.ExportFlags.Name,
		Fields:                   fields,
		//Filter: "EXISTS('{{ContactSegment[19014]}}')",
	}

	if *o.ExportFlags.MaxRecords > 0 {
		e.MaxRecords = *o.ExportFlags.MaxRecords
	}

	e, err = client.Contacts.CreateExport(ctx, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}
