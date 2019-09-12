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
//	"io"
	"errors"
	"strconv"
	"time"
	"strings"
	"regexp"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloquactl/pkg/printers"
	//cmdutil "github.com/elqx/eloquactl/pkg/util"
)

const (
	batchSize = 25000
	apiVersion = "2.0"
)

type Fields map[string]string

func NewCmdExport() *cobra.Command {
	cmd := &cobra.Command{
		Use: "export",
		Short: "export a resource from Eloqua",
		Long: ``,
		Example: "examples",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("You must specify the type of resource to export. See 'eloquactl export -h' for help and examples.")
		},
	}

	// create subcommands
	cmd.AddCommand(NewCmdExportActivities())
	cmd.AddCommand(NewCmdExportAccounts())
	cmd.AddCommand(NewCmdExportCdos())
	cmd.AddCommand(NewCmdExportContacts())

	return cmd
}

// export data given export definition
func export(ctx context.Context, ex *bulk.Export, keys *[]string, printer *printers.ResourcePrinter, client *bulk.BulkClient) {
	 // create sync definition
	sync, err := client.Syncs.Create(ctx, &bulk.Sync{SyncedInstanceURI: ex.Uri})
	if err != nil {
		fmt.Println(err)
	}

	// check sync status and download
	if err := waitSyncAndDownload(ctx, sync, keys,  printer, client); err != nil {
		 fmt.Println(err)
	}

}

func waitSyncAndDownload(ctx context.Context, sync *bulk.Sync, keys *[]string, printer *printers.ResourcePrinter, client *bulk.BulkClient) (error) {
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

	if err := download(ctx, syncId, keys, printer, client); err != nil {
		return err
	}

	return nil
}

func download(ctx context.Context, syncId int, keys *[]string, printer *printers.ResourcePrinter, client *bulk.BulkClient) (error) {
	opt := &bulk.QueryOptions{Limit: batchSize, Offset: 0}
	w := printers.NewTabWriter(os.Stdout)

	for {
		data, err := client.Syncs.GetData(ctx, syncId, opt)
		if err != nil {
			return err
		}

		(*printer).PrintResource(data.Items, w)
		w.Flush()

		if !data.HasMore {
			break
		}

		opt.Offset += batchSize
	}

	return nil
}

// parseFieldsStr parses fields string into a map of a field aliases and EML field representaions
// returns a slice of keys
func parseFieldsStr(str string, m *Fields) ([]string, error) {
	//m := make(map[string]string)
	s := strings.Split(str, ",")
	var k []string

	// looping over the slice and parsing its items si
	for _, si := range s {
		ss := strings.Split(si, ":")

		if len(ss) != 2 {
			return nil,  errors.New("Failed parsing fields string.")
		}

		k = append(k, strings.Trim(ss[0], " "))

		ss[0] = strings.Trim(ss[0], " ")
		ss[1] = strings.Trim(ss[1], " ")
		(*m)[ss[0]] = ss[1]
	}

	return k, nil
}

func checkDate(s string) error {
	re := regexp.MustCompile(DATE_REGEX)
	if match := re.MatchString(s); !match {
		return errors.New("invalid date string")
	}
	return nil
}

func checkISO8601(s string) error {
	// implementation missing
	return nil
}


func generateName() string {
	// implementation missing
	return "generated name"
}

func initClient() *bulk.BulkClient {
	auth := viper.GetStringMap("auth")
	bulkURL :=  strings.Replace(viper.GetString("bulkUrl"),"{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", auth["company"], auth["username"])
	password := auth["password"]

	tr := bulk.BasicAuthTransport{Username: username, Password: password.(string)}
	client := bulk.NewClient(bulkURL, tr.Client())

	return client
}
