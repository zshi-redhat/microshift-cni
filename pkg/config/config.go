package config

import (
	mcfg "github.com/openshift/microshift/pkg/config"
)

type MicroshiftConfig struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
	Config     *mcfg.MicroshiftConfig
}

func NewMicroshiftConfig() *MicroshiftConfig {
	return &MicroshiftConfig{
		Config: &mcfg.MicroshiftConfig{
			Cluster: mcfg.ClusterConfig{
				ClusterCIDR: "10.42.0.0/16",
				ServiceCIDR: "10.43.0.1/16",
			},
		},
	}
}
