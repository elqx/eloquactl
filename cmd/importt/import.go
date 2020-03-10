package importt

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ImportOptions struct {
}

func NewCmdImport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports data into Eloqua",
		Long:  "Imports data into Eloqua",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("executing import...")
		},
	}

	return cmd
}

func importt() {
	// create import definition
	// upload data (contacts, cdo data, accounts etc)
	// create sync
}
