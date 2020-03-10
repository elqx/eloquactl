package export

import (
	"context"
	"fmt"
	"os"
	"strconv"
	//"encoding/json"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
)

var (
	exportCdosLong = templates.LongDesc(`
		Export Eloqua CDO to a file or stdout.

		JSON and CSV file formats are supported`)

	exportCdosExample = templates.Examples(`
		# Export CDO given its name
		eloquactl export cdo NAME`)
)

type ExportCdosOptions struct {
	Client      func() *bulk.BulkClient
	PrintFlags  *cmdutil.PrintFlags
	ExportFlags *cmdutil.ExportFlags
}

func NewExportCdosOptions() *ExportCdosOptions {
	return &ExportCdosOptions{
		Client:      initClient,
		ExportFlags: cmdutil.NewExportFlags(),
		PrintFlags:  cmdutil.NewPrintFlags(),
	}
}

func NewCmdExportCdos() *cobra.Command {
	o := NewExportCdosOptions()

	cmd := &cobra.Command{
		Use:     "cdos <NAME>",
		Aliases: []string{"cdo"},
		Short:   "Export CDO to a file or stdout.",
		Long:    exportCdosLong,
		Example: exportCdosExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.Er("export cdos needs a cdo name or id for the command")
			}

			o.Complete(cmd)
			o.Validate(cmd)
			o.Run(cmd, args)
		},
	}
	// Add shared flags
	o.ExportFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	return cmd
}

func (p *ExportCdosOptions) Complete(cmd *cobra.Command) error {
	/*
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
	*/
	return nil
}

func (p *ExportCdosOptions) Validate(cmd *cobra.Command) error {
	// validate shared flags
	if err := p.ExportFlags.Validate(); err != nil {
		return err
	}

	if err := p.PrintFlags.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *ExportCdosOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := p.Client()

	var parentId int
	// check if args[0] is numeric or string
	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		// name is given, find id
		// TODO: check cache first
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
			}
		}
	}

	fieldsStr, _ := cmd.Flags().GetString("fields")
	fields := Fields{}
	var keys []string
	if fieldsStr == "" {
		// getting default fields
		// TODO: default fields should be cached
		r, err := client.Cdos.ListFields(ctx, parentId)
		if err != nil {
			fmt.Println("Failed getting cdo fields definitions")
			os.Exit(1)
		}

		for _, f := range r.Items {
			fields[f.InternalName] = f.Statement
		}

	} else {
		k, err := parseFieldsStr(fieldsStr, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys = k
	}

	printer, err := p.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	// should have Filter struct in the client library
	//var filter strings.Builder
	e := &bulk.Export{
		AreSystemTimestampsInUTC: *p.ExportFlags.AreSystemTimestampsInUTC,
		AutoDeleteDuration:       *p.ExportFlags.StagingFlags.AutoDeleteDuration,
		DataRetentionDuration:    *p.ExportFlags.StagingFlags.DataRetentionDuration,
		Name:                     *p.ExportFlags.Name,
		Fields:                   fields,
		// Filter: "'{{Contact.CreatedAt}}' >= '2019-09-01'",
		//MaxRecords: *p.ExportFlags.MaxRecords,
	}

	if *p.ExportFlags.MaxRecords > 0 {
		e.MaxRecords = *p.ExportFlags.MaxRecords
	}

	e, err = client.Cdos.CreateExport(ctx, parentId, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}
