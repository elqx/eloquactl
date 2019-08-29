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
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdExport() *cobra.Command {
	// exportCmd represents the export command
	exportCmd := &cobra.Command{
		Use: "export es",
		Short: "export a resource from Eloqua",
		Long: `long descritpion`,
		Example: "examples",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("export called")
		},
	}

	// create subcommands
	exportCmd.AddCommand(NewCmdExportActivities())
	exportCmd.AddCommand(NewCmdExportAccounts())
	exportCmd.AddCommand(NewCmdExportCdos())
	exportCmd.AddCommand(NewCmdExportContacts())

	return exportCmd
}


//func init() {
//	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}
