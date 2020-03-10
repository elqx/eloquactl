package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/elqx/eloquactl/pkg/printers"
	"github.com/spf13/cobra"
)

func Er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

type StagingFlags struct {
	AutoDeleteDuration    *string
	DataRetentionDuration *string
}

func NewStagingFlags() *StagingFlags {
	autoDeleteDuration := "PT12H"
	dataRetentionDuration := "PT12H"
	return &StagingFlags{
		AutoDeleteDuration:    &autoDeleteDuration,
		DataRetentionDuration: &dataRetentionDuration,
	}
}

func (f *StagingFlags) AddFlags(cmd *cobra.Command) {
	if f.AutoDeleteDuration != nil {
		cmd.Flags().String("auto-delete-duration", *f.AutoDeleteDuration, "Time until the definition will be deleted, expressed using the ISO-8601 standard.")
	}

	if f.DataRetentionDuration != nil {
		cmd.Flags().String("data-retention-duration", *f.DataRetentionDuration, "The length of time exported data should remain in the staging area., expressed using the ISO-8601 standard.")
	}
}

func (f *StagingFlags) Validate() error {
	return nil
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

func GetFlagUint(cmd *cobra.Command, flag string) uint {
	u, err := cmd.Flags().GetUint(flag)
	if err != nil {
		fmt.Printf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return u
}

type PrintFlags struct {
	NoHeaders    *bool
	OutputFormat *string
}

func NewPrintFlags() *PrintFlags {
	outputFormat := ""
	noHeaders := false
	return &PrintFlags{
		NoHeaders:    &noHeaders,
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

func (f *PrintFlags) ToPrinter() (printers.ResourcePrinter, error) {
	var printer printers.ResourcePrinter
	outputFormat := strings.ToLower(*f.OutputFormat) // table,json,jsonpath,ndj,csv,custom-columns
	switch outputFormat {
	case "json":
		printer = &printers.JsonPrinter{}
	case "ndj":
		printer = &printers.NdjPrinter{}
	default:
		printer = &printers.TablePrinter{}
	}
	return printer, nil
}

func (f *PrintFlags) Validate() error {
	// check if output format is one of the allowed values
	return nil
}

type ExportFlags struct {
	StagingFlags             *StagingFlags
	AreSystemTimestampsInUTC *bool
	Name                     *string
	Fields                   *string
	MaxRecords               *uint
	Filter                   *string
}

func NewExportFlags() *ExportFlags {
	stagingFlags := NewStagingFlags()
	areSystemTimestampsInUTC := true
	name := ""
	fields := ""
	filter := ""
	maxRecords := uint(0)
	return &ExportFlags{
		StagingFlags:             stagingFlags,
		AreSystemTimestampsInUTC: &areSystemTimestampsInUTC,
		Name:                     &name,
		Fields:                   &fields,
		Filter:                   &filter,
		MaxRecords:               &maxRecords,
	}
}

func (f *ExportFlags) AddFlags(cmd *cobra.Command) {
	f.StagingFlags.AddFlags(cmd)

	if f.AreSystemTimestampsInUTC != nil {
		cmd.Flags().BoolVarP(f.AreSystemTimestampsInUTC, "utc", "u", *f.AreSystemTimestampsInUTC, "Whether or not system timestamps will be exported in UTC.")
	}

	if f.Fields != nil {
		cmd.Flags().StringVar(f.Fields, "fields", *f.Fields, "List of fields to be included in the export operation.")
	}

	if f.Filter != nil {
		cmd.Flags().StringVar(f.Filter, "filter", *f.Filter, "The filter parameter uses Eloqua Markup Language to only return certain results.")
	}

	if f.MaxRecords != nil {
		cmd.Flags().UintVar(f.MaxRecords, "max-records", *f.MaxRecords, "The maximum amount of records.")
	}

	if f.Name != nil {
		cmd.Flags().StringVarP(f.Name, "name", "n", *f.Name, "The name of the export definition.")
	}
}

func (f *ExportFlags) Validate() error {
	err := f.StagingFlags.Validate()
	if err != nil {
		return err
	}

	if len(*f.Name) > 100 {
		// return error
	}
	return nil
}

type ImportFlags struct {
	StagingFlags *StagingFlags

	Fields                           *string
	IdentifierFieldName              *string
	IsSyncTriggeredOnImport          *bool
	IsUpdatingMultipleMatchedRecords *bool
	Name                             *string
	NullIdentifierFieldName          *string
	SyncActions                      *[]string
	UpdateRule                       *string
}

func NewImportFlags() *ImportFlags {
	stagingFlags := NewStagingFlags()
	fields := ""
	identifierFieldName := ""
	isSyncTriggeredOnImport := false
	isUpdatingMultipleMatchedRecords := false
	name := ""
	nullIdentifierFieldName := ""
	syncActions := []string{}
	updateRule := ""
	return &ImportFlags{
		StagingFlags:                     stagingFlags,
		Fields:                           &fields,
		IdentifierFieldName:              &identifierFieldName,
		IsSyncTriggeredOnImport:          &isSyncTriggeredOnImport,
		IsUpdatingMultipleMatchedRecords: &isUpdatingMultipleMatchedRecords,
		Name:                             &name,
		NullIdentifierFieldName:          &nullIdentifierFieldName,
		SyncActions:                      &syncActions,
		UpdateRule:                       &updateRule,
	}
}

func (f *ImportFlags) AddFlags(cmd *cobra.Command) {

}

func (f *ImportFlags) Validate() error {
	return nil
}

/*
type validateStringFlagFn func(p *string) error
type validateIntFlagFn func(i *int) error

func AddStringFlagVar(cmd *cobra.Command, p *string,  name, shorthand, usage string, validate validateFn, opts ...flagOpts) {
	cmd.Flags().StringiVarP(p, name, shorthand, dflt, desc)
}
*/

type FileNameFlags struct {
	FileNames *[]string
	Recursive *bool
}

func NewFileNameFlags() *FileNameFlags {
	fileNames := []string{}
	recursive := false
	return &FileNameFlags{
		FileNames: &fileNames,
		Recursive: &recursive,
	}
}

func (f *FileNameFlags) AddFlags(cmd *cobra.Command) {
	if f == nil {
		return
	}

	if f.FileNames != nil {
		cmd.Flags().StringSliceVarP(f.FileNames, "filename", "f", *f.FileNames, "Filename, directory, or URL to files to use to create the resource")
	}

	if f.Recursive != nil {
		cmd.Flags().BoolVarP(f.Recursive, "recursive", "R", *f.Recursive, "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	}
}

func (f *FileNameFlags) Validate() error {
	return nil
}

type ListFlags struct {
	Count         *int
	Depth         *string
	LastUpdatedAt *string
	OrderBy       *string
	Page          *int
	Search        *string
}

func NewListFlags() *ListFlags {
	count := int(1000)
	depth := "minimal"
	lastUpdatedAt := ""
	orderBy := ""
	page := int(1)
	search := ""
	return &ListFlags{
		Count:         &count,
		Depth:         &depth,
		LastUpdatedAt: &lastUpdatedAt,
		OrderBy:       &orderBy,
		Page:          &page,
		Search:        &search,
	}
}

func (f *ListFlags) AddFlags(cmd *cobra.Command) {
	if f == nil {
		return
	}
	if f.Count != nil {
		cmd.Flags().IntVar(f.Count, "count", *f.Count, "Maximum number of entities to return. Must be less than or equal to 1000 and greater than or equal to 1.")
	}
	if f.Depth != nil {
		cmd.Flags().StringVar(f.Depth, "depth", *f.Depth, "Level of detail returned by the request. Eloqua APIs can retrieve entities at three different levels of depth: minimal, partial, and complete.")
	}
	if f.OrderBy != nil {
		cmd.Flags().StringVar(f.OrderBy, "sort-by", *f.OrderBy, "Specifies the field by which list results are ordered.")
	}
	if f.Page != nil {
		cmd.Flags().IntVar(f.Page, "page", *f.Page, "Specifies which page of entities to return (the count parameter defines the number of entities per page).")
	}
	if f.Search != nil {
		cmd.Flags().StringVar(f.Search, "filter", *f.Search, "Specifies the search criteria used to retrieve entities.")
	}
}
