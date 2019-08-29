package export

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdExportContacts() *cobra.Command {
	cmd := &cobra.Command{
		Use: "contacts --filter=....",
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exporting contacts")
		},
	}

	cmd.Flags().StringP("filter", "f", "", "Contact filter. EML statement")

	return cmd
}
