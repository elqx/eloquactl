/*
Copyright © 2019 elqx <ignotas.petrulis@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package export

import (
	"context"
	"fmt"
	"io"
	"errors"
	"strconv"
	"time"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/elqx/eloqua-go/eloqua/bulk"
)

const (
	batchSize = 25000
	apiVersion = "2.0"
)

type Fields map[string]string

var client *bulk.BulkClient

type ExportOptions struct {
	// Used in CDO export
	ParentId int

	// Bulk export definition
	Export *bulk.Export
}

type ExportFunc func(context.Context, *ExportOptions) (*bulk.Export, error)

type ExportFuncMap map[string]ExportFunc

var efm = ExportFuncMap{}

func (e ExportFuncMap) RegisterFunc(fKey string, f ExportFunc) error {
	if _, exists := e[fKey]; exists {
		return errors.New("Function already registered")
	}
	e[fKey] = f
	return nil
}

func (e ExportFuncMap) Execute(fKey string, ctx context.Context, opt *ExportOptions) (*bulk.Export, error) {
	if f, exists := e[fKey]; exists {
		ex, err := f(ctx, opt)
		if err != nil {
			return nil, err
		}

		return ex, nil
	}
	return nil, errors.New("export function does not exist")
}

func NewCmdExport() *cobra.Command {
	// cmd represents the export command
	cmd := &cobra.Command{
		Use: "export",
		Short: "export a resource from Eloqua",
		Long: ``,
		Example: "examples",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("You must specify the type of resource to export. See 'eloquactl export -h' for help and examples.")
		},
	}

	cmd.PersistentFlags().BoolP("utc", "u", true, "Whether or not system timestamps will be exported in UTC.")
	cmd.PersistentFlags().String("auto-delete-duration", "PT12H", "Time until the definition will be deleted, expressed using the ISO-8601 standard.")
	cmd.PersistentFlags().String("data-retention-duration", "PT12H", "The length of time exported data should remain in the staging area., expressed using the ISO-8601 standard.")
	cmd.PersistentFlags().StringP("name", "n", "", "The name of the export definition.")
	cmd.PersistentFlags().String("fields", "", "List of fields to be included in the export operation.")
	cmd.PersistentFlags().String("filter", "", "The filter parameter uses Eloqua Markup Language to only return certain results.")
	cmd.PersistentFlags().Int("max-records", 1000, "The maximum amount of records.")

	auth := viper.GetStringMap("auth")
	bulkURL :=  strings.Replace(viper.GetString("bulkUrl"),"{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", auth["company"], auth["username"])
	password := auth["password"]

	tr := bulk.BasicAuthTransport{Username: username, Password: password.(string)}
	client = bulk.NewClient(bulkURL, tr.Client())
	// create subcommands
	cmd.AddCommand(NewCmdExportActivities())
	cmd.AddCommand(NewCmdExportAccounts())
	cmd.AddCommand(NewCmdExportCdos())
	cmd.AddCommand(NewCmdExportContacts())

	return cmd
}

func export(fKey string, ctx context.Context, opt *ExportOptions, out io.Writer) {
	// create export definition
	ex, err := efm.Execute(fKey, ctx, opt)
	if err != nil {
		fmt.Println(err)
	}
	 // create sync definition
	sync, err := client.Syncs.Create(ctx, &bulk.Sync{SyncedInstanceURI: ex.Uri})
	if err != nil {
		fmt.Println(err)
	}

	// check sync status and download
	if err := waitSyncAndDownload(ctx, sync, out); err != nil {
		 fmt.Println(err)
	}
}

func waitSyncAndDownload(ctx context.Context, sync *bulk.Sync, out io.Writer) (error) {
	syncId, err := strconv.Atoi(sync.Uri[7:])
	if err != nil {
		return err
	}

	for sync.Status != "success" && sync.Status != "error" {
		time.Sleep(2 * time.Second)
		sync, err = client.Syncs.Get(ctx, syncId)
		if err != nil {
			return errors.New("Failed to check sync status")
		}
	}

	if sync.Status == "error" {
		return errors.New("Failed to sync")
	}

	if err := download(ctx, syncId, out); err != nil {
		return err
	}

	return nil
}

func download(ctx context.Context, syncId int, out io.Writer) (error) {
	opt := &bulk.QueryOptions{Limit: batchSize, Offset: 0}

	for {
		data, err := client.Syncs.GetData(ctx, syncId, opt)
		if err != nil {
			return err
		}
		// TODO: format print
		for _, item := range data.Items {
			var str strings.Builder

			for _, v := range item {
				str.WriteString(v + "\t")
			}
			io.WriteString(out, str.String() + "\n")
		}

		if !data.HasMore {
			break
		}

		opt.Offset += batchSize
	}

	return nil
}

// parseFieldsStr parses fields string into a map of a field aliases and EML field representaions
func parseFieldsStr(str string, m *Fields) (error) {
	//m := make(map[string]string)
	s := strings.Split(str, ",")

	// looping over the slice and parsing its items si
	for _, si := range s {
		ss := strings.Split(si, ":")

		if len(ss) != 2 {
			return errors.New("Failed parsing fields string.")
		}

		ss[0] = strings.Trim(ss[0], " ")
		ss[1] = strings.Trim(ss[1], " ")
		(*m)[ss[0]] = ss[1]
	}

	return  nil
}
