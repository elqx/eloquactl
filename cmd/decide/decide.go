package decide

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdDecide() *cobra.Command {
	// cmd represents the decide command
	cmd := &cobra.Command{
		Use:     "decide --decision=...",
		Short:   "import contacts to the decision service",
		Long:    `long description`,
		Example: "examples",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("decide called")
		},
	}

	cmd.Flags().StringP("instance", "i", "", "Decision service instance id")

	return cmd
}
