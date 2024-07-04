package prodiverregistry

import (
	"context"
	"errors"

	sysv1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/constants"
	clientset "bytetrade.io/web3os/system-server/pkg/generated/clientset/versioned"
	"bytetrade.io/web3os/system-server/pkg/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var ErrProviderNotFound = errors.New("provider not found")

type Registry struct {
	registryClientset clientset.Interface
	namespace         string
}

func NewRegistry(clientset clientset.Interface) *Registry {
	registry := &Registry{
		registryClientset: clientset,
		namespace:         constants.MyNamespace,
	}

	return registry
}

func (r *Registry) GetProvider(ctx context.Context, dataType, group, version string) (*sysv1alpha1.ProviderRegistry, error) {
	providerRegistries, err := r.registryClientset.SysV1alpha1().
		ProviderRegistries(r.namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(providerRegistries.Items) > 0 {

		for _, pr := range providerRegistries.Items {
			if pr.Status.State == sysv1alpha1.Active {
				if pr.Spec.DataType == dataType &&
					pr.Spec.Group == group &&
					pr.Spec.Version == version &&
					pr.Spec.Kind == sysv1alpha1.Provider {
					return &pr, nil
				}
			}
		}

	}

	return nil, ErrProviderNotFound
}

func (r *Registry) GetWatchers(ctx context.Context, dataType, group, version string) ([]*sysv1alpha1.ProviderRegistry, error) {
	providerRegistries, err := r.registryClientset.SysV1alpha1().
		ProviderRegistries(r.namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	prs := make([]*sysv1alpha1.ProviderRegistry, 0)
	if len(providerRegistries.Items) > 0 {

		for _, pr := range providerRegistries.Items {
			if pr.Status.State == sysv1alpha1.Active {
				if pr.Spec.DataType == dataType &&
					pr.Spec.Group == group &&
					pr.Spec.Version == version &&
					pr.Spec.Kind == sysv1alpha1.Watcher {
					klog.Info("watcher callbacks, ", utils.PrettyJSON(pr))

					prs = append(prs, pr.DeepCopy())
				}
			}
		}

	}

	return prs, nil
}
