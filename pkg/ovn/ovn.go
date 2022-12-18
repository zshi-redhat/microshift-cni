package ovn

import (
	"path/filepath"

	mcfg "github.com/openshift/microshift/pkg/config"
	"github.com/spf13/pflag"
	"github.com/zshi-redhat/microshift-cni/pkg/assets"
	"github.com/zshi-redhat/microshift-cni/pkg/config"
	"github.com/zshi-redhat/microshift-cni/pkg/utils"
	"k8s.io/klog/v2"
)

const (
	manifestDir = "/etc/microshift/ovn"
)

func InstallOVNKubernetes(shiftConfig *config.MicroshiftConfig, flags *pflag.FlagSet) error {
	var err error

	err = ReadFromCmdLine(shiftConfig, flags)
	if err != nil {
		return err
	}

	ovnConfig, err := newOVNKubernetesConfig()
	if err != nil {
		return err
	}

	err = createGatewayBridges()
	if err != nil {
		return err
	}

	err = applyManifests(ovnConfig, shiftConfig)

	return nil
}

func ReadFromCmdLine(shiftConfig *config.MicroshiftConfig, flags *pflag.FlagSet) error {
	if s, err := flags.GetString("kubeconfig"); err == nil && flags.Changed("kubeconfig") {
		shiftConfig.Kubeconfig = s
	}
	if s, err := flags.GetString("cluster-cidr"); err == nil && flags.Changed("cluster-cidr") {
		shiftConfig.Config.Cluster.ClusterCIDR = s
	}
	if s, err := flags.GetString("service-cidr"); err == nil && flags.Changed("service-cidr") {
		shiftConfig.Config.Cluster.ServiceCIDR = s
	}
	return nil
}

// createGatewayBridges creates gateway bridges: br-ex and br-ex1
func createGatewayBridges() error {
	// TODO: use node IP from k8s node object to configure br-ex
	return utils.RunCommand("configure-ovs.sh", "OVNKubernetes")

}

func applyManifests(ovnConfig *OVNKubernetesConfig, shiftConfig *config.MicroshiftConfig) error {
	var (
		ns = []string{
			"ovn-kubernetes/namespace.yaml",
		}
		sa = []string{
			"ovn-kubernetes/node/serviceaccount.yaml",
			"ovn-kubernetes/master/serviceaccount.yaml",
		}
		r = []string{
			"ovn-kubernetes/role.yaml",
		}
		rb = []string{
			"ovn-kubernetes/rolebinding.yaml",
		}
		cr = []string{
			"ovn-kubernetes/clusterrole.yaml",
		}
		crb = []string{
			"ovn-kubernetes/clusterrolebinding.yaml",
		}
		cm = []string{
			"ovn-kubernetes/configmap.yaml",
		}
		apps = []string{
			"ovn-kubernetes/master/daemonset.yaml",
			"ovn-kubernetes/node/daemonset.yaml",
		}
	)

	var kubeconfigPath string
	if shiftConfig.Kubeconfig == "" {
		kubeconfigPath = shiftConfig.Config.KubeConfigPath(mcfg.KubeAdmin)
	}

	if err := assets.ApplyNamespaces(ns, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply ns %v: %v", ns, err)
		return err
	}
	if err := assets.ApplyServiceAccounts(sa, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply serviceAccount %v %v", sa, err)
		return err
	}
	if err := assets.ApplyRoles(r, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply role %v: %v", r, err)
		return err
	}
	if err := assets.ApplyRoleBindings(rb, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply rolebinding %v: %v", rb, err)
		return err
	}
	if err := assets.ApplyClusterRoles(cr, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply clusterRole %v %v", cr, err)
		return err
	}
	if err := assets.ApplyClusterRoleBindings(crb, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply clusterRoleBinding %v %v", crb, err)
		return err
	}
	extraParams := assets.RenderParams{
		"MTU":                       ovnConfig.MTU,
		"KubeconfigPath":            kubeconfigPath,
		"KubeconfigDir":             filepath.Dir(kubeconfigPath),
		"ClusterCIDR":               shiftConfig.Config.Cluster.ClusterCIDR,
		"ServiceCIDR":               shiftConfig.Config.Cluster.ServiceCIDR,
		"ovn_kubernetes_microshift": "quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:34faa03671a950607aec06a2ba515426d09f80a38101d2e91a508b5fd5316116",
	}
	if err := assets.ApplyConfigMaps(cm, extraParams, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply configMap %v %v", cm, err)
		return err
	}
	if err := assets.ApplyDaemonSets(apps, extraParams, kubeconfigPath); err != nil {
		klog.Warningf("Failed to apply apps %v %v", apps, err)
		return err
	}

	/*
		manifests, err := assets.List()
		if err != nil {
			return err
		}

		for _, m := range manifests {
			fmt.Println(f)
		}

		for _, m := range manifests {
			bytes, err := assets.Asset(m)
			if err != nil {
				return
			}
			fmt.Println(string(bytes))
		}
	*/

	/*
		objs := []*uns.Unstructured{}

		data := render.MakeRenderData()
		data.Data["MTU"] = ovnConfig.MTU

		manifests, err := render.RenderDirs(manifestDir, &data)
		if err != nil {
			return err
		}
		objs = append(objs, manifests...)

		for _, obj := range objs {
			apply.ApplyObject(obj)
		}
	*/
	return nil
}
