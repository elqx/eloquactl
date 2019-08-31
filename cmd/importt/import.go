package importt

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdImport() *cobra.Command {
	cmd := &cobra.Command{
		Use: "import",
		Short: "Imports data into Eloqua",
		Long: "Imports data into Eloqua",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("executing import...")
		},
	}

	return cmd
}

