package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addRunFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
}

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

func Manifest() {
	out := make([][]byte, 0)

	files, err := embedded.List()
	if err != nil {
		return
	}

	for _, f := range files {
		bytes, err := embedded.Asset(f)
		if err != nil {
			return
		}
		out = append(out, bytes)
	}
	fmt.Println(out)
}
