package export

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdExportCdos() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cdos <NAME>",
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exporting cdos")
		},
	}

	return cmd
}
