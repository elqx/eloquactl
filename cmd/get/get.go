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
package get

import (
	//	"context"
	"fmt"
	//"errors"
	//"strconv"
	//"time"
	"strings"
	//"regexp"
	//"os"

	"github.com/elqx/eloqua-go/eloqua/bulk"
	"github.com/elqx/eloqua-go/eloqua/pkg/auth"
	"github.com/elqx/eloqua-go/eloqua/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"github.com/elqx/eloquactl/pkg/printers"
	//cmdutil "github.com/elqx/eloquactl/pkg/util"
)

const (
	apiVersion = "2.0"
)

// type Fields map[string]string
func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get a resource from Eloqua",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Getting...")
			fmt.Println("Done.")
		},
	}

	// create subcommands
	cmd.AddCommand(NewCmdGetCampaigns())
	cmd.AddCommand(NewCmdGetCdoFields())
	cmd.AddCommand(NewCmdGetContactFields())
	cmd.AddCommand(NewCmdGetEmails())
	cmd.AddCommand(NewCmdGetEmailGroups())
	cmd.AddCommand(NewCmdGetForms())

	return cmd
}

func initClient() *bulk.BulkClient {
	bauth := viper.GetStringMap("auth")
	bulkURL := strings.Replace(viper.GetString("bulkUrl"), "{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", bauth["company"], bauth["username"])
	password := bauth["password"]

	tr := auth.BasicAuthTransport{Username: username, Password: password.(string)}
	client := bulk.NewClient(bulkURL, tr.Client())

	return client
}

func initRestClient() *rest.RestClient {
	bauth := viper.GetStringMap("auth")
	restURL := strings.Replace(viper.GetString("restUrl"), "{version}", apiVersion, 1)
	username := fmt.Sprintf("%v\\%v", bauth["company"], bauth["username"])
	password := bauth["password"]

	tr := auth.BasicAuthTransport{Username: username, Password: password.(string)}
	client := rest.NewClient(restURL, tr.Client())

	return client
}
