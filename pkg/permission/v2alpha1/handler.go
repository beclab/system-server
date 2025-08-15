package v2alpha1

import (
	"context"
	"errors"
	"fmt"

	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api/response"
	"bytetrade.io/web3os/system-server/pkg/constants"
	"bytetrade.io/web3os/system-server/pkg/utils"
	"github.com/emicklei/go-restful/v3"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

type handler struct {
	kubeClient kubernetes.Interface
}

func (h *handler) register(req *restful.Request, resp *restful.Response) {
	ok, username := h.validate(req, resp)
	if !ok {
		return
	}

	var app string
	if ok, app = h.execute(req, resp, username, h.bindingProvider); !ok {
		return
	}

	klog.Info("success to register provider, ", username, ", app=", app)
	reg := &RegisterResp{}
	response.Success(resp, reg)
}

func (h *handler) unregister(req *restful.Request, resp *restful.Response) {
	ok, username := h.validate(req, resp)
	if !ok {
		return
	}

	_, app := h.execute(req, resp, username, h.unbindingProvider)

	klog.Info("success to unregister provider, ", username, ", app=", app)
	response.SuccessNoData(resp)
}

func (h *handler) execute(req *restful.Request, resp *restful.Response, username string,
	action func(ctx context.Context, user, app, serviceAccount string, roles []*rbacv1.ClusterRole) error) (success bool, appName string) {
	var err error
	var perm PermissionRegister

	if err = req.ReadEntity(&perm); err != nil {
		api.HandleError(resp, req, err)
		return
	}

	if perm.App == "" {
		err = errors.New("invalid app, app name is empty")
		klog.Error(err)
		api.HandleError(resp, req, err)
		return
	}

	appName = perm.App

	for _, p := range perm.Perm {
		if p.ProviderName == "" {
			continue
		}

		if p.ServiceAccount == nil {
			p.ServiceAccount = ptr.To("default")
		}

		providerName := h.userProviderRef(username, p.ProviderName)
		roles, err := h.getProvider(req.Request.Context(), username, providerName)
		if err != nil {
			klog.Error("fail to get provider roles, ", err)
			api.HandleError(resp, req, err)
			return
		}

		if len(roles) == 0 {
			klog.Warning("no roles found for provider, ", providerName)
			continue
		}

		if err = action(req.Request.Context(), username, perm.App, *p.ServiceAccount, roles); err != nil {
			klog.Error("fail to bind provider, ", err)
			api.HandleError(resp, req, err)
			return
		}

	} // end of for app perm loop

	success = true
	return
}

func (h *handler) validate(req *restful.Request, resp *restful.Response) (ok bool, username string) {
	account, err := h.getUser(req, resp)
	if err != nil {
		klog.Error(err)
		api.HandleUnauthorized(resp, req, err)
		return
	}

	var namespace string
	// get provider of user
	if isSA, saNamespace, _ := utils.IsServiceAccount(account); isSA {
		var isUserNamespace bool
		if isUserNamespace, username = utils.IsUserNamespace(saNamespace); !isUserNamespace {
			err := fmt.Errorf("user is not found in namespace %s", saNamespace)
			klog.Error(err)
			api.HandleUnauthorized(resp, req, err)
			return
		} // end of service account check

		namespace = saNamespace
	} else {
		username = account
		namespace = "user-system-" + username
	}

	klog.Infof("User %s is a system user in namespace %s", username, namespace)
	if constants.MyNamespace != namespace && constants.MyUserspace != namespace {
		api.HandleUnauthorized(resp, req, fmt.Errorf("invalid user, %s", username))
		return
	}

	ok = true
	return
}
