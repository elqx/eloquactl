package feed

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed --instance=...",
		Short: "Feed contacts into a campaign using the specified feeder instance.",
		Long:  "Feed contacts into a campaign using the specified feeder instance.",
		Example: `  # Feed contacts into a campaign from the feeder instance a12d53dd1
  eloquactl feed --instance=a12d53dd1
		`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("feeding...")
		},
	}

	return cmd
}
