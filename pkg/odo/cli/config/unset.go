package config

import (
	"fmt"
	"strings"

	"github.com/redhat-developer/odo/pkg/odo/cli/ui"

	"github.com/pkg/errors"
	"github.com/redhat-developer/odo/pkg/config"
	"github.com/redhat-developer/odo/pkg/odo/util"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubernetes/pkg/kubectl/cmd/templates"
)

const unsetCommandName = "unset"

var (
	unsetLongDesc = ktemplates.LongDesc(`Unset an individual value in the Odo configuration file.

%[1]s
%[2]s
`)
	unsetExample = ktemplates.Examples(`
   # Unset a configuration value in the local config
   %[1]s %[2]s
   %[1]s %[3]s 
   %[1]s %[4]s  
   %[1]s %[5]s 
   %[1]s %[6]s 
   %[1]s %[7]s  
   %[1]s %[8]s 
   %[1]s %[9]s 
   %[1]s %[10]s
	`)
)

// UnsetOptions encapsulates the options for the command
type UnsetOptions struct {
	paramName       string
	configForceFlag bool
}

// NewUnsetOptions creates a new UnsetOptions instance
func NewUnsetOptions() *UnsetOptions {
	return &UnsetOptions{}
}

// Complete completes UnsetOptions after they've been created
func (o *UnsetOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	o.paramName = args[0]
	return
}

// Validate validates the UnsetOptions based on completed values
func (o *UnsetOptions) Validate() (err error) {
	return
}

// Run contains the logic for the command
func (o *UnsetOptions) Run() (err error) {

	cfg, err := config.New()

	if err != nil {
		return errors.Wrapf(err, "")
	}

	if value, ok := cfg.GetConfiguration(o.paramName); ok && (value != nil) {
		if !o.configForceFlag && !ui.Proceed(fmt.Sprintf("Do you want to unset %s in the config", o.paramName)) {
			fmt.Println("Aborted by the user.")
			return nil
		}
		err = cfg.DeleteConfiguration(strings.ToLower(o.paramName))
		if err != nil {
			return err
		}

		fmt.Println("Local config was successfully updated.")

		return nil
		// if its found but nil then show the error
	} else if ok && (value == nil) {
		return errors.New("config already unset, cannot unset a configuration which is not set")
	}
	return errors.New(o.paramName + " is not a valid configuration variable")

}

// NewCmdUnset implements the config unset odo command
func NewCmdUnset(name, fullName string) *cobra.Command {
	o := NewUnsetOptions()
	configurationUnsetCmd := &cobra.Command{
		Use:   name,
		Short: "Unset a value in odo config file",
		Long:  fmt.Sprintf(unsetLongDesc, config.FormatLocallySupportedParameters()),
		Example: fmt.Sprintf(fmt.Sprint("\n", unsetExample), fullName,
			config.ComponentType, config.ComponentName, config.MinMemory, config.MaxMemory, config.Memory, config.Ignore, config.MinCPU, config.MaxCPU, config.CPU),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("please provide a parameter name")
			} else if len(args) > 1 {
				return fmt.Errorf("only one parameter is allowed")
			} else {
				return nil
			}

		}, Run: func(cmd *cobra.Command, args []string) {
			util.LogErrorAndExit(o.Complete(name, cmd, args), "")
			util.LogErrorAndExit(o.Validate(), "")
			util.LogErrorAndExit(o.Run(), "")
		},
	}
	configurationUnsetCmd.Flags().BoolVarP(&o.configForceFlag, "force", "f", false, "Don't ask for confirmation, unsetting the config directly")

	return configurationUnsetCmd
}
