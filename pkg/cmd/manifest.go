package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
	embedded "github.com/zshi-redhat/microshift-cni/assets"
)

func addRunFlags(cmd *cobra.Command) {}

func NewManifestCommand() *cobra.Command {
	//	cfg := config.newOVNKubernetesConfigFromFile()

	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Get OVNKubernetes Manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Manifest()
		},
	}
	addRunFlags(cmd)
	return cmd
}

func Manifest() error {
	files, err := embedded.List()
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f)
	}

	/* print manifest content to stdout

	for _, f := range files {
		bytes, err := embedded.Asset(f)
		if err != nil {
			return
		}
		fmt.Println(string(bytes))
	}

	*/
	return nil
}
