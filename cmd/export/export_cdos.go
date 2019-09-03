package export

import (
	"fmt"
	"strconv"
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/util"
)

func NewCmdExportCdos() *cobra.Command {
	fKey := "cdo"
	cmd := &cobra.Command{
		Use: "cdos <NAME>",
		Aliases: []string{"cdo"},
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				util.Er("export cdos needs a cdo name for the command")
			}

			ctx := context.Background()
			var parentId int
			// check if args[0] is numeric or string
			parentId , err := strconv.Atoi(args[0])
			if err != nil {
				// name is given, find id
				// TODO: check cache first
				cdos, err := client.Cdos.List(ctx)
				if err != nil {
					util.Er("Failed listing CDOs")
				}

				for _, cdo := range cdos.Items {
					if cdo.Name == args[0] {
						// skipping /customObjects/ prefix
						parentId, err = strconv.Atoi(cdo.Uri[15:])
						if err != nil {
							util.Er("Error extracting custom object id")
						}
					}
				}
			}
			// find cdo id for a name

			maxRecords, _ := cmd.Flags().GetInt("max-records")
			areSystemTimestampsInUTC, _ := cmd.Flags().GetBool("utc")
			autoDeleteDuration, _ := cmd.Flags().GetString("auto-delete-duration")
			dataRetentionDuration, _ := cmd.Flags().GetString("data-retention-duration")
			defaultName := "cdo export"

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
				AreSystemTimestampsInUTC: areSystemTimestampsInUTC,
				AutoDeleteDuration: autoDeleteDuration,
				DataRetentionDuration: dataRetentionDuration,
				Name: defaultName,
				Fields: fields,
				// Filter: "'{{Contact.CreatedAt}}' >= '2019-09-01'",
				MaxRecords: maxRecords,
			}
			// exporting activities
			opt := &ExportOptions{ParentId: parentId, Export: e}
			export(fKey, ctx, opt, os.Stdout)
		},
	}
	efm.RegisterFunc(fKey, func(ctx context.Context, opt *ExportOptions) (*bulk.Export, error) {
		return client.Cdos.CreateExport(ctx, opt.ParentId, opt.Export)
	})
	return cmd
}
