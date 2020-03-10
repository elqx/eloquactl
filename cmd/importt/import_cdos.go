package importt

import (
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/spf13/cobra"
)

type ImportCdoOptions struct {
}

func NewCmdImportCdosCommand() {
	cmd := &cobra.Command{
		Use:     "cdos",
		Aliases: []string{"cdo"},
		Short:   "",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmdutil.AddStagingFlags(cmd)
	cmdutil.AddMapDataCardsFlags(cmd)
	cmdutil.AddNameFlag(cmd)
}
