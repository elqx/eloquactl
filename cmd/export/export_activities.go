package export

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/elqx/eloqua-go/eloqua/bulk"
)

type Fields map[string]string

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
)

func NewCmdExportActivities() *cobra.Command {
	fKey := "activities"

	cmd := &cobra.Command{
		Use: "activities --type=type --format=...",
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
			if fieldsStr == "" {
				fields = activityFields[tflag]
			} else {
				err := parseFieldsStr(fieldsStr, &fields)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			defaultFilter := fmt.Sprintf("'{{Activity.Type}}' = '%v' AND '{{Activity.CreatedAt}}' >= '2019-08-30'", t)
			e := &bulk.Export{
				AreSystemTimestampsInUTC: areSystemTimestampsInUTC,
				AutoDeleteDuration: autoDeleteDuration,
				DataRetentionDuration: dataRetentionDuration,
				Name: defaultName,
				Fields: fields,
				Filter: defaultFilter,
				MaxRecords: maxRecords,
			}
			// exporting activities
			export(fKey, e, os.Stdout)
			/*
			auth := viper.GetStringMap("auth")
			username := fmt.Sprintf("%v\\%v", auth["company"].(string), auth["username"].(string))
			password := auth["password"].(string)
			accountInfo, err := bulk.GetAccountInfo(username, password)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(accountInfo.Urls.Apis.Rest.Bulk)
			*/
		},
	}

	cmd.Flags().StringP("type", "t", "", "Activity type")
	// required flags
	cmd.MarkFlagRequired("type")
	//cmd.Flags().StringP("format", "f", "CSV", "Data format. Possible values: CSV, JSON. Default value: CSV.")
	// register activities export function
	efm.RegisterFunc(fKey, client.Activities.CreateExport)
	return cmd
}
