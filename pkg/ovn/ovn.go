package ovn

import (
	"path/filepath"

	"github.com/zshi-redhat/microshift-cni/pkg/assets"
	"github.com/zshi-redhat/microshift-cni/pkg/config"
	"github.com/zshi-redhat/microshift-cni/pkg/utils"
	"k8s.io/klog/v2"
)

const (
	manifestDir = "/etc/microshift/ovn"
)

func InstallOVNKubernetes(shiftConfig *config.MicroshiftConfig) error {
	var err error

	cfg, err := newOVNKubernetesConfig()
	if err != nil {
		return err
	}

	err = createGatewayBridges()
	if err != nil {
		return err
	}

	err = applyManifests(cfg, shiftConfig)

	return nil
}

// createGatewayBridges creates gateway bridges: br-ex and br-ex1
func createGatewayBridges() error {
	return utils.RunCommand("configure-ovs.sh", "OVNKubernetes")

}

func applyManifests(cfg *OVNKubernetesConfig, shiftCfg *config.MicroshiftConfig) error {
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

	kubeconfigPath := shiftCfg.Kubeconfig

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
		"MTU":            cfg.MTU,
		"KubeconfigPath": kubeconfigPath,
		"KubeconfigDir":  filepath.Join("/var/lib", "/resources/kubeadmin"),
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
		data.Data["MTU"] = cfg.MTU

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
