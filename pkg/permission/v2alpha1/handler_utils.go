package v2alpha1

import (
	"context"
	"fmt"
	"net/http"

	"bytetrade.io/web3os/system-server/pkg/constants"
	providerv2alpha1 "bytetrade.io/web3os/system-server/pkg/providerregistry/v2alpha1"
	"github.com/emicklei/go-restful/v3"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (h *handler) getUser(req *restful.Request, resp *restful.Response) (string, error) {
	user := req.Request.Header.Get(constants.BflUserKey)
	if user == "" {
		err := restful.NewError(http.StatusUnauthorized, "User not found in request header")
		klog.Error(err)
		return "", err
	}

	return user, nil
}

func (h *handler) getProvider(ctx context.Context, user string, providerName string) ([]*rbacv1.ClusterRole, error) {
	clusterRoles, err := h.kubeClient.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Errorf("Failed to list cluster roles: %v", err)
		return nil, err
	}

	var roles []*rbacv1.ClusterRole
	for _, role := range clusterRoles.Items {
		if providerRef, ok := role.Annotations[providerv2alpha1.ProviderRefAnnotation]; ok {
			if providerRef == h.userProviderRef(user, providerName) {
				klog.Infof("Found provider role: %s for user: %s", role.Name, user)
				roles = append(roles, &role)
			}
		}
	}

	return roles, nil
}

func (h *handler) bindingProvider(ctx context.Context, user, app, serviceAccount string, roles []*rbacv1.ClusterRole) error {
	appNamespace := fmt.Sprintf("%s-%s", app, user)
	for _, role := range roles {
		binding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: h.getProviderBindingName(appNamespace, serviceAccount, role.Name),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      serviceAccount,
					Namespace: appNamespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.SchemeGroupVersion.Group,
				Kind:     role.Kind,
				Name:     role.Name,
			},
		}

		if rb, err := h.kubeClient.RbacV1().ClusterRoleBindings().Get(ctx, binding.Name, metav1.GetOptions{}); err == nil {
			klog.Infof("Cluster role binding %s already exists for service account %s in app %s", rb.Name, serviceAccount, appNamespace)
			rb.Subjects = binding.Subjects
			rb.RoleRef = binding.RoleRef
			if _, err := h.kubeClient.RbacV1().ClusterRoleBindings().Update(ctx, rb, metav1.UpdateOptions{}); err != nil {
				klog.Errorf("Failed to update cluster role binding %s: %v", rb.Name, err)
				return err
			}
			continue
		}

		if _, err := h.kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, binding, metav1.CreateOptions{}); err != nil {
			klog.Errorf("Failed to create cluster role binding %s: %v", binding.Name, err)
			return err
		}
		klog.Infof("Created cluster role binding %s for service account %s in app %s", binding.Name, serviceAccount, appNamespace)
	}
	return nil
}

func (h *handler) unbindingProvider(ctx context.Context, user, app, serviceAccount string, roles []*rbacv1.ClusterRole) error {
	appNamespace := fmt.Sprintf("%s-%s", app, user)
	for _, role := range roles {
		bindingName := h.getProviderBindingName(appNamespace, serviceAccount, role.Name)
		if err := h.kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, bindingName, metav1.DeleteOptions{}); err == nil {
			klog.Infof("Deleted cluster role binding %s for service account %s in app %s", bindingName, serviceAccount, appNamespace)
		} else {
			klog.Errorf("Failed to delete cluster role binding %s for service account %s in app %s: %v", bindingName, serviceAccount, appNamespace, err)
		}
	}

	return nil
}

func (h *handler) userProviderRef(user, providerName string) string {
	return fmt.Sprintf("user-system-%s/%s", user, providerName)
}

func (h *handler) getProviderBindingName(appNamespace, serviceAccount, roleName string) string {
	return fmt.Sprintf("%s:%s:%s", appNamespace, serviceAccount, roleName)
}
