package export

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/elqx/eloqua-go/eloqua/bulk"
)


var (
	activityTypes = map[string]string{
		"es": "EmailSend",
		"eo": "EmailOpen",
		"ec": "EmailClickthrough",
		"bb": "Bounceback",
		"fs": "FormSubmit",
		"su": "Subscribe",
		"un": "Unsubscribe",
		"wv": "WebVisit",
		"pv": "PageView",
	}
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

func NewCmdExportActivities() *cobra.Command {
	fKey := "activities"

	cmd := &cobra.Command{
		Use: "activities --type=type --format=...",
		Aliases: []string{"activity"},
		Short: "",
		Long: "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			tflag, _ := cmd.Flags().GetString("type")
			t, exists := activityTypes[tflag]
			if !exists {
				// should print error, help and exit
				fmt.Printf("--type value %v is not supported", tflag)
			}
			maxRecords, _ := cmd.Flags().GetInt("max-records")
			areSystemTimestampsInUTC, _ := cmd.Flags().GetBool("utc")
			autoDeleteDuration, _ := cmd.Flags().GetString("auto-delete-duration")
			dataRetentionDuration, _ := cmd.Flags().GetString("data-retention-duration")
			defaultName := fmt.Sprintf("activities %v export", t)

			fieldsStr, _ := cmd.Flags().GetString("fields")

			fields := Fields{}
			ctx := context.Background()

			if fieldsStr == "" {
				// getting default activity fields for the activity type
				opt := &bulk.ActivityFieldListQueryOptions{ActivityType: t}
				// TODO: default fields should be cached 
				r, err := client.Activities.ListFields(ctx, opt)
				if err != nil {
					fmt.Printf("Failed to list activity fields for activity type: %v", t)
					os.Exit(1)
				}

				for _, f := range r.Items {
					fields[f.InternalName] = f.Statement
				}
			} else {
				err := parseFieldsStr(fieldsStr, &fields)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			since, _ := cmd.Flags().GetString("since")
			until, _ := cmd.Flags().GetString("until")

			// should have Filter struct in the client library
			var filter strings.Builder
			filter.WriteString(fmt.Sprintf("'{{Activity.Type}}' = '%v'", t))
			if since != "" && until == "" {
				filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' >= '%v'", since))
			}

			if since == "" && until != "" {
				filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' < '%v'", until))
			}

			if since != "" && until != "" {
				filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' >= '%v'", since))
				filter.WriteString(fmt.Sprintf(" AND '{{Activity.CreatedAt}}' < '%v'", until))
			}

			e := &bulk.Export{
				AreSystemTimestampsInUTC: areSystemTimestampsInUTC,
				AutoDeleteDuration: autoDeleteDuration,
				DataRetentionDuration: dataRetentionDuration,
				Name: defaultName,
				Fields: fields,
				Filter: filter.String(),
				MaxRecords: maxRecords,
			}
			// exporting activities
			export(fKey, ctx, e, os.Stdout)
		},
	}
	cmd.Flags().StringP("type", "t", "", "Activity type")
	cmd.Flags().String("since", "", "The lower bound of the date range filter (inclusive).")
	cmd.Flags().String("until", "", "The upper bound of the date range filter (noninclusive).")
	// required flags
	cmd.MarkFlagRequired("type")
	//cmd.Flags().StringP("format", "f", "CSV", "Data format. Possible values: CSV, JSON. Default value: CSV.")
	// register activities export function
	efm.RegisterFunc(fKey, client.Activities.CreateExport)
	return cmd
}
