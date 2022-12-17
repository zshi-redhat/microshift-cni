package config

type MicroshiftConfig struct {
	// absolute path to Kubernetes kubeconfig file
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

func NewMicroshiftConfig() {
	return &MicroshiftConfig{}
}
