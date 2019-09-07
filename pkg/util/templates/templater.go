package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type templater struct {
	HelpTemplate string
	UsageTemplate string
	CommandGroups
}

func (t *templater) HelpFunc() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, s []string) {
		tpl := template.New("help")
		tpl.Funcs(t.templateFuncs())
		template.Must(tpl.Parse(t.HelpTemplate))
		err := tpl.Execute(os.Stdout, cmd)
		if err != nil {
			cmd.Println(err)
		}

	}
}

func (t *templater) UsageFunc() func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		tpl := template.New("usage")
		tpl.Funcs(t.templateFuncs())
		template.Must(tpl.Parse(t.UsageTemplate))
		return tpl.Execute(os.Stdout, cmd)
	}
}
/*
func (t *templater) cmdGroups(cmd *cobra.Command, all []*cobra.Command) []CommandGroup {
}
*/
func (t *templater) cmdGroupsString(cmd *cobra.Command) string {
	groups := []string{}
	for _, cmdGroup := range t.CommandGroups {
		cmds := []string{cmdGroup.Message}
		for _, c := range cmdGroup.Commands {
			cmds = append(cmds, "  " + rpad(c.Name(), c.NamePadding()) + " " + c.Short)
		}
		groups = append(groups, strings.Join(cmds, "\n"))
	}
	return strings.Join(groups, "\n\n")
}

func (t *templater) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"trim": strings.TrimSpace,
		"cmdGroupsString": t.cmdGroupsString,
	}
}

func ActsAsRootCommand(cmd *cobra.Command, groups ...CommandGroup) {
	if cmd == nil {
		panic("nil root command")
	}
	templater := &templater{
		UsageTemplate: UsageTemplate(),
		HelpTemplate: HelpTemplate(),
		CommandGroups: groups,
	}
	cmd.SetHelpFunc(templater.HelpFunc())
	cmd.SetUsageFunc(templater.UsageFunc())
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}
