package export

/*
import (
	"github.com/spf13/cobra"
)

type CompleteFn func()
type ValidateFn func() error

type Command struct {
	*cobra.Command

	completeFns []*CompleteFn
	validateFns []*ValidateFn
}

func (c *Command) CompleteFlags() {
	// run all completeFns
	for _, complete := range c.completeFns {
		complete()
	}
}

func (c *Command) ValidateFlags() error {
	// run all validateFns
	for _, validate := range c.validateFns {

		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) AddCommand(cmd *cobra.Command) {
}


func AddFlagString(cmd *Command, name, shorthand string, complete *CompleteFn, validate *ValidateFn) {
}

type CmdRunner func(*CmdConfig) error

type CmdConfig struct {
}

func CmdBuilder(parent *Command, cr CmdRunner, name, shortDesc, longDesc, example string, options ...cmdOption) *cobra.Command {
	cmd := &cobra.Command{
		Use: name,
		Short: shortDesc,
		Long: longDesc,
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cr(c))
		},
	}

	if parent != nil {
		parent.AddCommand(cmd)
	}
}
*/
