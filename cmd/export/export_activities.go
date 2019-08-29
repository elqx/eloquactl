package export

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdExportActivities() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activities --type=type --format=...",
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			t, _ := cmd.Flags().GetString("type")
			if t == "" {
				fmt.Println("exporting activities")
			} else {
				fmt.Printf("exporting %v activities\n", t)
			}
		},
	}

	cmd.Flags().StringP("type", "t", "", "Activity type")
	//cmd.Flags().StringP("format", "f", "CSV", "Data format. Possible values: CSV, JSON. Default value: CSV.")

	return cmd
}
