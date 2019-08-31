package export

import (
	"fmt"
	//"strings"
	"os"

	"github.com/spf13/cobra"
	"github.com/elqx/eloqua-go/eloqua/bulk"
)

func NewCmdExportContacts() *cobra.Command {
	// should be a const to prevent modification?
	fKey := "contacts"

	cmd := &cobra.Command{
		Use: "contacts --filter=....",
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exporting contacts")
			// should parse flags
			// initialize export definition
			defaultName := "contact export"
			defaultFields := make(map[string]string)
			defaultFields["EmailAddress"] = "{{Contact.Field(C_EmailAddress)}}"
			defaultFilter := "'{{Contact.CreatedAt}} >= '2019-08-30''"

			e := &bulk.Export{
				Name: defaultName,
				Fields: defaultFields , // should be a map of strings
				Filter: defaultFilter, // should be created based on parsed params
			}
			export(fKey, e, os.Stdout)
		},
	}

	cmd.Flags().StringP("filter", "f", "", "Contact filter. EML statement")
	efm.RegisterFunc(fKey, client.Contacts.CreateExport)
	return cmd
}
