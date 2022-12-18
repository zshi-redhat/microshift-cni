package config

import (
	mcfg "github.com/openshift/microshift/pkg/config"
)

func NewMicroshiftConfig() *mcfg.MicroshiftConfig {
	return &mcfg.MicroshiftConfig{
		Cluster: mcfg.ClusterConfig{
			ClusterCIDR: "10.42.0.0/16",
			ServiceCIDR: "10.43.0.1/16",
		},
	}
}
