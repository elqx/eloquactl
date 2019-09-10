package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/elqx/eloquactl/pkg/printers"
)

func Er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func AddStagingFlags(cmd *cobra.Command) {
	cmd.Flags().String("auto-delete-duration", "PT12H", "Time until the definition will be deleted, expressed using the ISO-8601 standard.")
	cmd.Flags().String("data-retention-duration", "PT12H", "The length of time exported data should remain in the staging area., expressed using the ISO-8601 standard.")
}

func AddNameFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("name", "n", "", "The name of the export definition.")
}

func AddMapDataCardsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("map-data-cards", false, "Whether or not custom object records or event registrants will be mapped on import. If you set it to true, you must specify the fields for mapping.")
	cmd.Flags().Bool("map-data-cards-case-sensitive-match", false, "Whether to perform a case sensitive search when mapping custom object records or events to a contact or account.")
	cmd.Flags().String("map-data-cards-entity-fields", "", "Specifies which Eloqua entity field will be used for mapping.")
	cmd.Flags().String("map-data-cards-entity-type", "", "Specifies the entity of the custom object record or event import. Allowed values are 'Contact' or 'Company'.")
	cmd.Flags().String("map-data-cards-source-field", "", "Specifies the source field that will be used for matching.")
}

func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		fmt.Printf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

func GetFlagStringSlice(cmd *cobra.Command, flag string) []string {
	ss, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		fmt.Printf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return ss
}

func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		fmt.Printf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return b
}

func GetFlagInt(cmd *cobra.Command, flag string) int {
	i, err := cmd.Flags().GetInt(flag)
	if err != nil {
		fmt.Printf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

type PrintFlags struct {
	NoHeaders *bool
	OutputFormat *string
}

func NewPrintFlags() *PrintFlags {
	outputFormat := ""
	noHeaders := false
	return &PrintFlags{
		NoHeaders: &noHeaders,
		OutputFormat: &outputFormat,
	}
}

func (f *PrintFlags) AddFlags(cmd *cobra.Command) {
	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(f.OutputFormat, "output", "o", *f.OutputFormat, "Output format. One of: json|table|name|jsonpath.")
	}

	if f.NoHeaders != nil {
		cmd.Flags().BoolVar(f.NoHeaders, "no-headers", *f.NoHeaders, "When using the default (table) output format, don't print headers (default print headers).")
	}
}

func (f * PrintFlags) ToPrinter() (printers.ResourcePrinter, error) {
	var printer printers.ResourcePrinter
	outputFormat := strings.ToLower(*f.OutputFormat)
	fmt.Println("FORMAT", outputFormat)
	switch outputFormat {
		case "json":
			printer = &printers.JsonPrinter{}
		default:
			printer = &printers.TablePrinter{}
	}
	return printer, nil
}
