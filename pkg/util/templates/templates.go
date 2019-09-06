package templates

import (
	"strings"
)

const (
	SectionAliases = `{{if gt (len .Aliases) 0}}Aliases:
{{.NameAndAliases}}

{{end}}`

	SectionExamples = `{{if .HasExample}}Examples:
{{trim .Example}}

{{end}}`

	SectionSubcommands = `{{if .HasAvailableSubCommands}}{{cmdGroupsString .}}

{{end}}`
)

// template for help command
func HelpTemplate() string {
	return "{{with or .Long .Short}}{{ . | trim}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}"
}

// command usage template
func UsageTemplate() string {
	sections := []string{
		"\n\n",
		SectionAliases,
		SectionExamples,
		SectionSubcommands,
	}
	return strings.TrimSpace(strings.Join(sections, ""))
}
