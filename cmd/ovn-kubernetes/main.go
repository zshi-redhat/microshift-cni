package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zshi-redhat/microshift-cni/pkg/config"
	ovn "github.com/zshi-redhat/microshift-cni/pkg/ovn-kubernetes"

	"k8s.io/component-base/cli"
)

func main() {
	command := newCommand()
	code := cli.Run(command)
	os.Exit(code)
}

func addFlags(cmd *cobra.Command, cfg *config.MicroShiftConfig) {
	flags := cmd.Flags()
	cmd.Flags().String("kubeconfig", cfg.Kubeconfig, "The absolute path for Kubernetes kubeconfig file")
}

func newCommand() *cobra.Command {
	cfg := config.NewMicroshiftConfig()

	cmd := &cobra.Command{
		Use:   "ovn-kubernetes",
		Short: "MicroShift OVNKubernetes CNI",
	}
	addFlags(cmd, cfg)

	cmd.Run = func(cmd *cobra.Command, args []string) error {
		return ovn.InstallOVNKubernetes(cfg)
	}
	return cmd
}
