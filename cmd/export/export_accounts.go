package export

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdExportAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts",
		Short:   "",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exporting accounts")
		},
	}

	return cmd
}
