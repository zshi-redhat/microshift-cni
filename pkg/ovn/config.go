package ovn

import (
	"errors"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

const (
	configFile = "/etc/microshift/ovn.yaml"
)

type OVNKubernetesConfig struct {
	// Configuration for microshift-ovs-init.service
	OVSInit OVSInit `json:"ovsInit,omitempty"`
	// MTU to use for the geneve tunnel interface.
	// This must be 100 bytes smaller than the uplink mtu.
	// Default is 1400.
	MTU uint32 `json:"mtu,omitempty"`
}

type OVSInit struct {
	// disable microshift-ovs-init.service.
	// OVS bridge "br-ex" needs to be configured manually when disableOVSInit is true.
	DisableOVSInit bool `json:"disableOVSInit,omitempty"`
	// Uplink interface for OVS bridge "br-ex"
	GatewayInterface string `json:"gatewayInterface,omitempty"`
	// Uplink interface for OVS bridge "br-ex1"
	ExternalGatewayInterface string `json:"externalGatewayInterface,omitempty"`
}

func (o *OVNKubernetesConfig) ValidateOVSBridge(bridge string) error {
	_, err := net.InterfaceByName(bridge)
	if err != nil {
		return err
	}
	return nil
}

func (o *OVNKubernetesConfig) withDefaults() *OVNKubernetesConfig {
	o.OVSInit.DisableOVSInit = false
	o.MTU = 1400
	return o
}

func newOVNKubernetesConfigFromFile(path string) (*OVNKubernetesConfig, error) {
	o := new(OVNKubernetesConfig)
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &o)
	if err != nil {
		return nil, fmt.Errorf("parsing OVNKubernetes config: %v", err)
	}
	return o, nil
}

func newOVNKubernetesConfig() (*OVNKubernetesConfig, error) {
	if _, err := os.Stat(configFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			klog.Info("OVNKubernetes config file not found, assuming default values")
			return new(OVNKubernetesConfig).withDefaults(), nil
		}
		return nil, fmt.Errorf("failed to get OVNKubernetes config file: %v", err)
	}

	o, err := newOVNKubernetesConfigFromFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("getting OVNKubernetes config: %v", err)
	}

	klog.Info("got OVNKubernetes config from file %q", configFile)
	return o, nil
}
