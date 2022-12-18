package main

import (
	"os"

	mcfg "github.com/openshift/microshift/pkg/config"
	"github.com/spf13/cobra"
	"github.com/zshi-redhat/microshift-cni/pkg/config"
	"github.com/zshi-redhat/microshift-cni/pkg/ovn"

	"k8s.io/component-base/cli"
)

func main() {
	command := newCommand()
	code := cli.Run(command)
	os.Exit(code)
}

func addFlags(cmd *cobra.Command, cfg *mcfg.MicroshiftConfig) {
	cmd.Flags().String("servicecidr", cfg.Cluster.ServiceCIDR, "The Kubernetes service network range")
	cmd.Flags().String("clustercidr", cfg.Cluster.ClusterCIDR, "The Kubernetes cluster network range")
}

func newCommand() *cobra.Command {
	shiftConfig := config.NewMicroshiftConfig()

	cmd := &cobra.Command{
		Use:   "ovn-kubernetes",
		Short: "MicroShift OVNKubernetes CNI",
	}
	addFlags(cmd, shiftConfig)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return ovn.InstallOVNKubernetes(shiftConfig)
	}
	return cmd
}
