package main

import (
	"os"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/component-base/cli"

	cmds "github.com/zshi-redhat/microshift-cni/pkg/cmd"
)

func main() {
	command := newCommand()
	code := cli.Run(command)
	os.Exit(code)
}

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "OVNKubernetes",
		Short: "MicroShift OVNKubernetes plugin",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	originalHelpFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		originalHelpFunc(command, strings)
	})

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	cmd.AddCommand(cmds.NewManifestCommand())
	return cmd
}
