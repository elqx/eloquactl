package export

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
)

const (
	DATE_REGEX = "\\d{4}-\\d{2}-\\d{2}"
)

var (
	activityTypes = map[string]bool{
		"EmailSend":         true,
		"EmailOpen":         true,
		"EmailClickthrough": true,
		"Bounceback":        true,
		"FormSubmit":        true,
		"Subscribe":         true,
		"Unsubscribe":       true,
		"WebVisit":          true,
		"PageView":          true,
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

type Validator interface {
	Validate() error
}

type ValidatorFunc func() error

func (fn ValidatorFunc) Validate() error {
	return fn()
}

// ExportActivitiesOptions declare the arguments accepted by the 'export activities' command
// this struct should have all configurable properties of an export
type ExportActivitiesOptions struct {
	Client func() *bulk.BulkClient

	PrintFlags  *cmdutil.PrintFlags
	ExportFlags *cmdutil.ExportFlags

	// Command specific options.
	// Used to construct the filter.
	ActivityType string
	Since        string
	Until        string

	// inherits Validator method
	//Validate ValidatorFunc
}

func NewExportActivitiesOptions() *ExportActivitiesOptions {
	return &ExportActivitiesOptions{
		Client:      initClient,
		ExportFlags: cmdutil.NewExportFlags(),
		PrintFlags:  cmdutil.NewPrintFlags(),

		//Validate: AggErrorValidator,
	}
}

func NewCmdExportActivities() *cobra.Command {
	o := NewExportActivitiesOptions()

	cmd := &cobra.Command{
		Use:     "activities --type ACTIVITYTYPE",
		Aliases: []string{"activity"},
		Short:   "Export Eloqua activities to a file or stdout.",
		Long:    exportActivitiesLong,
		Example: exportActivitiesExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Validate()
			o.Run(cmd)
		},
	}
	// Add shared flags
	o.ExportFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	// Add flags specific to activities export
	//AddStringFlag(cmd, o.ActivityType, "type", "t", "", "Activity type", checkDate, requiredOpt)
	cmd.Flags().StringVarP(&o.ActivityType, "type", "t", "", "Activity type")
	cmd.Flags().StringVar(&o.Since, "since", "", "The lower bound of the date range filter (inclusive).")
	cmd.Flags().StringVar(&o.Until, "until", "", "The upper bound of the date range filter (noninclusive).")

	// Required flags
	cmd.MarkFlagRequired("type")

	return cmd
}

// Complete completes the options provided
func (p *ExportActivitiesOptions) Complete(cmd *cobra.Command) {
	// StagingFlags, ExportFlags and PrintFlags are completed
	// here should only be the completion of the filter option
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

	str := filter.String()
	p.ExportFlags.Filter = &str
}

// Validate validates the options provided
func (p *ExportActivitiesOptions) Validate() error {
	//var e Errors
	// validate shared flags
	if err := p.ExportFlags.Validate(); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := p.PrintFlags.Validate(); err != nil {
		//e = aggregateError(err)
		return err
	}

	// validate command specific flags
	if _, exists := activityTypes[p.ActivityType]; !exists {
		// should print error, help and exit
		//err := errors.New("Unsuported activity type")
		//e = aggregateError(err)
		return errors.New("No such activity type")
	}

	if err := checkDate(p.Since); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := checkDate(p.Until); err != nil {
		//e = aggregateError(err)
		return err
	}

	return nil
}

// Run executes the command
func (p *ExportActivitiesOptions) Run(cmd *cobra.Command) error {
	ctx := context.Background()
	client := p.Client()
	// if fields are empty, should get the fields via api
	// fields is a runtime option if not provided
	fields := Fields{}
	var keys []string
	if len(*p.ExportFlags.Fields) == 0 {
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
		k, err := parseFieldsStr(*p.ExportFlags.Fields, &fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys = k
	}

	if len(*p.ExportFlags.Filter) == 0 {
		// get fields via api and construct the filter
		// fields should be cached
	}

	printer, err := p.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	e := &bulk.Export{
		AreSystemTimestampsInUTC: *p.ExportFlags.AreSystemTimestampsInUTC,
		AutoDeleteDuration:       *p.ExportFlags.StagingFlags.AutoDeleteDuration,
		DataRetentionDuration:    *p.ExportFlags.StagingFlags.DataRetentionDuration,
		Name:                     *p.ExportFlags.Name,
		Fields:                   fields,
		Filter:                   *p.ExportFlags.Filter,
		//	MaxRecords: p.MaxRecords,
	}

	if *p.ExportFlags.MaxRecords > 0 {
		e.MaxRecords = *p.ExportFlags.MaxRecords
	}

	data, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return err
	}

	data = append(data, '\n')

	e, err = client.Activities.CreateExport(ctx, e)
	if err != nil {
		return err
	}

	export(ctx, e, &keys, &printer, client)

	return nil
}
