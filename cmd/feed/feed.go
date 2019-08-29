package feed

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use: "feed --instance=...",
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("feeding...")
		},
	}

	return cmd
}
