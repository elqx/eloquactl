package importt

import (
	//"context"
	cmdutil "github.com/elqx/eloquactl/pkg/util"
	"github.com/elqx/eloquactl/pkg/util/templates"
	"github.com/spf13/cobra"
	//"github.com/elqx/eloqua-go/eloqua/bulk"
)

var (
	importContactsLong = templates.LongDesc(`
		Import contacts to Eloqua from a file or stdin.
		
		JSON and CSV file formats are supported`)

	importContactsExample = templates.Examples(`
		# Import contacts into Eloqua from a file
		eloquactl import contacts -f=contacts.csv`)
)

type ImportContactsOptions struct {
	ImportFlags   *cmdutil.ImportFlags
	FileNameFlags *cmdutil.FileNameFlags
}

func NewImportContactsOptions() *ImportContactsOptions {
	return &ImportContactsOptions{
		ImportFlags:   cmdutil.NewImportFlags(),
		FileNameFlags: cmdutil.NewFileNameFlags(),
	}
}

func NewCmdImportContacts() *cobra.Command {
	o := NewImportContactsOptions()
	cmd := &cobra.Command{
		Use:     "contacts",
		Aliases: []string{"contact"},
		Short:   "Export Eloqua contacts to a file or stdout",
		Long:    importContactsLong,
		Example: importContactsExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.Run(cmd, args)
		},
	}

	// Add shared falgs
	o.ImportFlags.AddFlags(cmd)
	o.FileNameFlags.AddFlags(cmd)

	return cmd
}

func (o *ImportContactsOptions) Validate() error {
	//var e Errors
	if err := o.ImportFlags.Validate(); err != nil {
		//e = aggregateError(err)
		return err
	}

	if err := o.FileNameFlags.Validate(); err != nil {
		//e = aggregateError(err)
		return err
	}

	return nil
}

func (o *ImportContactsOptions) Run(cmd *cobra.Command, args []string) {
	/*
		ctx := context.Background()
		i := &bulk.Import{
			AutoDeleteDuration: o.ImportFlags.AutoDeleteDuration,
			DataRetentionDuration: o.ImportFlags.DataRetentionDuration,
			Fields: o.ImportFlags.Fields, //fix
			IdentifierFieldName: o.ImportFlags.Identifier.FieldName,
			IsSyncTriggeredOnImport: o.ImportFlags.IsSyncTriggeredOnImport,
			IsUpdatingMultipleMatchedRecords: o.ImportFlags.IsUpdatingMultipleMatchedRecords,
			Name: o.ImportFlags.Name,
			NullIdentifierFieldName: o.ImportFlags.NullIdentifierName,
			SyncActions: o.ImportFlags.SyncActions, //fix
			UpdateRule: o.ImportFlags.UpdateRule,
		}

		i, err = client.Contacts.CreateImport(ctx, &i)
		if err != nil {
			return err
		}

		// read contacts from file
		client.Contacts.Upload()
	*/
}
