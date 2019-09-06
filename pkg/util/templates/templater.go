package templates

import (
	"github.com/spf13/cobra"
)

type templater struct {
	HelpTemplate string
	UsageTemplate string
}

func (t *templater) HelpFunc() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, s []string) {
		// use text/template package
	}
}

func (t *templater) UsageFunc() func(*cobra.Command) error{
	return func(cmd *cobra.Command) error {
		// use some template and exexute
		// text/template package
		return nil
	}
}

func ActsAsRootCommand(cmd *cobra.Command, groups ...CommandGroups) {
	if cmd == nil {
		panic("nil root command")
	}
	templater := &templater{
		UsageTemplate: UsageTemplate(),
		HelpTemplate: HelpTemplate(),
	}
	cmd.SetHelpFunc(templater.HelpFunc())
	cmd.SetUsageFunc(templater.UsageFunc())
}


