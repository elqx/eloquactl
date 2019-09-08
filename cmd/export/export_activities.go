package export

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/util/templates"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
)

const (
	DATE_REGEX = "\\d{4}-\\d{2}-\\d{2}"
	EXPORT_ACTIVITIES_KEY = "activities"
)

var (
	activityTypes = map[string]bool{
		"EmailSend": true,
		"EmailOpen": true,
		"EmailClickthrough": true,
		"Bounceback": true,
		"FormSubmit": true,
		"Subscribe": true,
		"Unsubscribe": true,
		"WebVisit": true,
		"PageView": true,
	}

	exportActivitiesLong = templates.LongDesc(`
		Export Eloqua activities to a file or stdout.

		JSON and CSV file formats are supported`)

	exportActivitiesExample = templates.Examples(`
		# Export EmailSend activities given date ranges
		eloquactl export activities --type=EmailSend --since=2019-01-01 --until=2019-02-01

		# Export specific fields of EmailOpen activities
		eloquactl export activities --type=EmailOpen --since=2019-01-01 --fields='ActivityDate:{{Activity.CreatedAt}},EmailAddress:{{Activity.Field(EmailAddress)}}'`)

/*
	activityFields = map[string]Fields{
		"es": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"SubjectLine": "{{Activity.Field(SubjectLine)}}",
			"EmailWebLink": "{{Activity.Field(EmailWebLink)}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"DeploymentId": "{{Activity.Field(EmailDeploymentId)}}",
			"EmailSendType": "{{Activity.Field(EmailSendType)}}",
			"CampaignResponseDate": "{{Activity.CampaignResponse.CreatedAt}}",
			"CampaignResponseMemberStatus": "{{Activity.CampaignResponse.Field(MemberStatus)}}",
		},
		"eo": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"IpAddress": "{{Activity.Field(IpAddress)}}",
			"VisitorId": "{{Activity.Visitor.Id}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"SubjectLine": "{{Activity.Field(SubjectLine)}}",
			"EmailWebLink": "{{Activity.Field(EmailWebLink)}}",
			"VisitorExternalId": "{{Activity.Visitor.ExternalId}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"DeploymentId": "{{Activity.Field(EmailDeploymentId)}}",
			"EmailSendType": "{{Activity.Field(EmailSendType)}}",
			"CampaignResponseDate": "{{Activity.CampaignResponse.CreatedAt}}",
			"CampaignResponseMemberStatus": "{{Activity.CampaignResponse.Field(MemberStatus)}}",
		},
		"ec": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"IpAddress": "{{Activity.Field(IpAddress)}}",
			"VisitorId": "{{Activity.Visitor.Id}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"SubjectLine": "{{Activity.Field(SubjectLine)}}",
			"EmailWebLink": "{{Activity.Field(EmailWebLink)}}",
			"EmailClickedThruLink": "{{Activity.Field(EmailClickedThruLink)}}",
			"VisitorExternalId": "{{Activity.Visitor.ExternalId}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"DeploymentId": "{{Activity.Field(EmailDeploymentId)}}",
			"EmailSendType": "{{Activity.Field(EmailSendType)}}",
			"CampaignResponseDate": "{{Activity.CampaignResponse.CreatedAt}}",
			"CampaignResponseMemberStatus": "{{Activity.CampaignResponse.Field(MemberStatus)}}",
		},
		"bb": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"DeploymentId": "{{Activity.Field(EmailDeploymentId)}}",
			"SmtpErrorCode": "{{Activity.Field(SmtpErrorCode)}}",
			"SmtpStatusCode": "{{Activity.Field(SmtpStatusCode)}}",
			"SmtpMessage": "{{Activity.Field(SmtpMessage)}}",

		},
		"fs": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"VisitorId": "{{Activity.Visitor.Id}}",
			"VisitorExternalId": "{{Activity.Visitor.ExternalId}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"RawData": "{{Activity.Field(RawData)}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"CampaignResponseDate": "{{Activity.CampaignResponse.CreatedAt}}",
			"CampaignResponseMemberStatus": "{{Activity.CampaignResponse.Field(MemberStatus)}}",
		},
		"su": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
		},
		"un": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"AssetId": "{{Activity.Asset.Id}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"EmailAddress": "{{Activity.Field(EmailAddress)}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"EmailRecipientId": "{{Activity.Field(EmailRecipientId)}}",
			"AssetType": "{{Activity.Asset.Type}}",
			"AssetName": "{{Activity.Asset.Name}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"ExternalId": "{{Activity.ExternalId}}",
		},
		"wv": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"VisitorId": "{{Activity.Visitor.Id}}",
			"VisitorExternalId": "{{Activity.Visitor.ExternalId}}",
			"ReferrerUrl": "{{Activity.Field(ReferrerUrl)}}",
			"IpAddress": "{{Activity.Field(IpAddress)}}",
			"NumberOfPages": "{{Activity.Field(NumberOfPages)}}",
			"FirstPageViewUrl": "{{Activity.Field(FirstPageViewUrl)}}",
			"Duration": "{{Activity.Field(Duration)}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"LinkedToContactDate": "{{Activity.Field(LinkedToContactDate)}}",
			"WebVisitSavedId": "{{Activity.Field(WebVisitSavedId)}}",
		},
		"pv": Fields{
			"ActivityId": "{{Activity.Id}}",
			"ActivityType": "{{Activity.Type}}",
			"ActivityDate": "{{Activity.CreatedAt}}",
			"ContactId": "{{Activity.Contact.Id}}",
			"CampaignId": "{{Activity.Campaign.Id}}",
			"VisitorId": "{{Activity.Visitor.Id}}",
			"VisitorExternalId": "{{Activity.Visitor.ExternalId}}",
			"WebVisitId": "{{Activity.Field(WebVisitId)}}",
			"Url": "{{Activity.Field(Url)}}",
			"ReferrerUrl": "{{Activity.Field(ReferrerUrl)}}",
			"IpAddress": "{{Activity.Field(IpAddress)}}",
			"IsWebTrackingOptedIn": "{{Activity.Field(IsWebTrackingOptedIn)}}",
			"ExternalId": "{{Activity.ExternalId}}",
			"LinkedToContactDate": "{{Activity.Field(LinkedToContactDate)}}",
			"PageViewSavedId": "{{Activity.Field(PageViewSavedId)}}",
			"CampaignResponseDate": "{{Activity.CampaignResponse.CreatedAt}}",
			"CampaignResponseMemberStatus": "{{Activity.CampaignResponse.Field(MemberStatus)}}",
		},

	}
	*/
)


// ExportActivitiesOptions declare the arguments accepted by the 'export activities' command
// this struct should have all configurable properties of an export
type ExportActivitiesOptions struct {
	UTC bool
	AutoDeleteDuration string
	DataRetentionDuration string
	Name string
	Fields string
	Filter string // raw filter should not be exposed to the user, but it can be in this struct (flags are exposed to the user, not this struct)
	MaxRecords int
	ActivityType string
	Since string
	Until string
}

func NewCmdExportActivities() *cobra.Command {
	o := &ExportActivitiesOptions{}

	cmd := &cobra.Command{
		Use: "activities --type ACTIVITYTYPE",
		Aliases: []string{"activity"},
		Short: "Export Eloqua activities to a file or stdout.",
		Long: exportActivitiesLong,
		Example: exportActivitiesExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd)
		},
	}
	// bind flags to options struct fields
	cmdutil.AddStagingFlags(cmd)
	cmdutil.AddNameFlag(cmd)
	cmd.Flags().BoolP("utc", "u", true, "Whether or not system timestamps will be exported in UTC.")
	cmd.Flags().String("fields", "", "List of fields to be included in the export operation.")
	//cmd.Flags().String("filter", "", "The filter parameter uses Eloqua Markup Language to only return certain results.")
	cmd.Flags().Int("max-records", 1000, "The maximum amount of records.")

	cmd.Flags().StringP("type", "t", "", "Activity type")
	cmd.Flags().String("since", "", "The lower bound of the date range filter (inclusive).")
	cmd.Flags().String("until", "", "The upper bound of the date range filter (noninclusive).")
	// required flags
	cmd.MarkFlagRequired("type")
	//cmd.Flags().StringP("format", "f", "CSV", "Data format. Possible values: CSV, JSON. Default value: CSV.")
	// register activities export function
	efm.RegisterFunc(EXPORT_ACTIVITIES_KEY, func(ctx context.Context, opt *ExportOptions) (*bulk.Export, error) {
		return client.Activities.CreateExport(ctx, opt.Export)
	})
	return cmd
}

// Complete completes the options with default values
func (p *ExportActivitiesOptions) Complete(cmd *cobra.Command) error {
	p.AutoDeleteDuration = cmdutil.GetFlagString(cmd, "auto-delete-duration")
	p.DataRetentionDuration = cmdutil.GetFlagString(cmd, "data-retention-duration")

	p.UTC = cmdutil.GetFlagBool(cmd, "utc")

	// name does not have a default, generate name if not specified
	p.Name = cmdutil.GetFlagString(cmd, "name")
	if len(p.Name) == 0 {
		p.Name = generateName()
	}
	// fields flag does not have a default value, use fields from api
	p.Fields = cmdutil.GetFlagString(cmd, "fields")

	p.MaxRecords = cmdutil.GetFlagInt(cmd, "max-records")
	p.Since = cmdutil.GetFlagString(cmd, "since")
	p.Until = cmdutil.GetFlagString(cmd, "until")
	p.ActivityType = cmdutil.GetFlagString(cmd, "type")

	// filter should constructed differently
	// Use harcoded fields?
	var filter strings.Builder
	filter.WriteString(fmt.Sprintf("'{{Activity.Type}}' = '%v'", p.ActivityType))
	if p.Since != "" && p.Until == "" {
		filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' >= '%v'", p.Since))
	}

	if p.Since == "" && p.Until != "" {
		filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' < '%v'", p.Until))
	}

	if p.Since != "" && p.Until != "" {
		filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' >= '%v'", p.Since))
		filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' < '%v'", p.Until))
	}

	p.Filter = filter.String()

	return nil
}

// Validate validates the options provided
func (p *ExportActivitiesOptions) Validate() error {
	// validate activity type
	if _, exists := activityTypes[p.ActivityType]; !exists {
		// should print error, help and exit
		fmt.Printf("--type value %v is not supported", p.ActivityType)
	}

	// TODO: check that staging options follow the ISO-8601 standard
	if err := checkISO8601(p.AutoDeleteDuration); err != nil {
		fmt.Println("Failed validating auto-delete-duration. Value should follow ISO-8601 standard.")
	}

	if err := checkISO8601(p.DataRetentionDuration); err != nil {
		fmt.Println("Failed validation data-retention-duration. Value should follow ISO-8601 standard.")
	}
	// check that name is no longer than 100 (based on doc)
	if len(p.Name) > 100 {
		fmt.Println("Export name must be no longer than 100 characters long.")
	}

	// validate fields if they were provided by the user
	if len(p.Fields) > 0 {
		// TODO: check regex expression
	}

	if p.MaxRecords < 0 {
		fmt.Println("max-records should be greater than zero")
	}

	if err := checkDate(p.Since); err != nil {
		return err
	}

	if err := checkDate(p.Since); err != nil {
		return err
	}

	return nil
}

// Run executes the command
func (p *ExportActivitiesOptions) Run(cmd *cobra.Command) error {
	ctx := context.Background()

	auth := viper.GetStringMap("auth")
	bulkURL :=  strings.Replace(viper.GetString("bulkUrl"),"{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", auth["company"], auth["username"])
	password := auth["password"]

	tr := &bulk.BasicAuthTransport{Username: username, Password: password.(string)}
	client = bulk.NewClient(bulkURL, tr.Client())
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	if len(p.Fields) == 0 {
		// get fields via api and assign
		opt := &bulk.ActivityFieldListQueryOptions{ActivityType: p.ActivityType}
		// TODO: default fields should be cached 
		r, err := client.Activities.ListFields(ctx, opt)
		if err != nil {
			fmt.Printf("Failed to list activity fields for activity type: %v", p.ActivityType)
			os.Exit(1)
		}

		for _, f := range r.Items {
			fields[f.InternalName] = f.Statement
		}
	} else {
		err := parseFieldsStr(p.Fields, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(p.Filter) == 0 {
		// get fields via api and construct the filter
		// fields should be cached
	}

	e := &bulk.Export{
		AreSystemTimestampsInUTC: p.UTC,
		AutoDeleteDuration: p.AutoDeleteDuration,
		DataRetentionDuration: p.DataRetentionDuration,
		Name: p.Name,
		Fields: fields,
		Filter: p.Filter,
		MaxRecords: p.MaxRecords,
	}

	opt := &ExportOptions{Export: e}
	export(EXPORT_ACTIVITIES_KEY, ctx, opt, os.Stdout)

	return nil
}
