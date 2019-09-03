package export

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/elqx/eloqua-go/eloqua/bulk"
)

func NewCmdExportContacts() *cobra.Command {
	// should be a const to prevent modification?
	fKey := "contacts"

	cmd := &cobra.Command{
		Use: "contacts --filter=....",
		Aliases: []string{"contact"},
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			maxRecords, _ := cmd.Flags().GetInt("max-records")
			areSystemTimestampsInUTC, _ := cmd.Flags().GetBool("utc")
			autoDeleteDuration, _ := cmd.Flags().GetString("auto-delete-duration")
			dataRetentionDuration, _ := cmd.Flags().GetString("data-retention-duration")
			defaultName := "contacts export"

			fieldsStr, _ := cmd.Flags().GetString("fields")
			ctx := context.Background()
			fields := Fields{}
			if fieldsStr == "" {
				// getting default fields
				// TODO: default fields should be cached
				r, err := client.Contacts.GetFields(ctx)
				if err != nil {
					fmt.Println("Failed getting contact fields definitions")
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
				Filter: "'{{Contact.CreatedAt}}' >= '2019-09-01'",
				MaxRecords: maxRecords,
			}
			// exporting activities
			opt := &ExportOptions{Export: e}
			export(fKey, ctx, opt, os.Stdout)
		},
	}
	//cmd.Flags().StringP("filter", "f", "", "Contact filter. EML statement")
	efm.RegisterFunc(fKey, func(ctx context.Context, opt *ExportOptions) (*bulk.Export, error) {
		return client.Contacts.CreateExport(ctx, opt.Export)
	})
	return cmd
}
