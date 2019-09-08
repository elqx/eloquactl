package export

import (
	"fmt"
	"strings"
	"strconv"
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/util/templates"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
)

const (
	EXPORT_CDOS_KEY = "cdos"
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
	UTC bool
	AutoDeleteDuration string
	DataRetentionDuration string
	Name string
	Fields string
	Filter string
	MaxRecords int
}

func NewCmdExportCdos() *cobra.Command {
	o := ExportCdosOptions{}
	cmd := &cobra.Command{
		Use: "cdos <NAME>",
		Aliases: []string{"cdo"},
		Short: "Export CDO to a file or stdout.",
		Long: exportCdosLong,
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
	cmdutil.AddStagingFlags(cmd)
	cmdutil.AddNameFlag(cmd)
	cmd.Flags().BoolP("utc", "u", true, "Whether or not system timestamps will be exported in UTC.")
	cmd.Flags().String("fields", "", "List of fields to be included in the export operation.")
	cmd.Flags().Int("max-records", 1000, "The maximum amount of records.")

	efm.RegisterFunc(EXPORT_CDOS_KEY, func(ctx context.Context, opt *ExportOptions) (*bulk.Export, error) {
		return client.Cdos.CreateExport(ctx, opt.ParentId, opt.Export)
	})

	return cmd
}

func (p *ExportCdosOptions) Complete(cmd *cobra.Command) error {
	p.AutoDeleteDuration = cmdutil.GetFlagString(cmd, "auto-delete-duration")
	p.DataRetentionDuration = cmdutil.GetFlagString(cmd, "data-retention-duration")

	p.UTC = cmdutil.GetFlagBool(cmd, "utc")

	// name does not have a default, generate name if not specified
	p.Name = cmdutil.GetFlagString(cmd, "name")
	if len(p.Name) == 0 {
		p.Name = generateName()
	}
	p.Fields = cmdutil.GetFlagString(cmd, "fields")

	p.MaxRecords = cmdutil.GetFlagInt(cmd, "max-records")

	return nil
}

func (p *ExportCdosOptions) Validate(cmd *cobra.Command) error {
	// TODO: check that staging options follow the ISO-8601 standard
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

	return nil
}

func (p *ExportCdosOptions) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	auth := viper.GetStringMap("auth")

	bulkURL :=  strings.Replace(viper.GetString("bulkUrl"),"{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", auth["company"], auth["username"])
	password := auth["password"]

	tr := &bulk.BasicAuthTransport{Username: username, Password: password.(string)}
	client = bulk.NewClient(bulkURL, tr.Client())

	var parentId int
	// check if args[0] is numeric or string
	parentId , err := strconv.Atoi(args[0])
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
		err := parseFieldsStr(fieldsStr, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// should have Filter struct in the client library
	//var filter strings.Builder
	e := &bulk.Export{
		AreSystemTimestampsInUTC: p.UTC,
		AutoDeleteDuration: p.AutoDeleteDuration,
		DataRetentionDuration: p.DataRetentionDuration,
		Name: p.Name,
		Fields: fields,
		// Filter: "'{{Contact.CreatedAt}}' >= '2019-09-01'",
		MaxRecords: p.MaxRecords,
	}

	opt := &ExportOptions{ParentId: parentId, Export: e}
	export(EXPORT_CDOS_KEY, ctx, opt, os.Stdout)
	return nil
}
