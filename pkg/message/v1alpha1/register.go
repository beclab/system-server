package message

import (
	"context"
	syserrs "errors"
	"flag"
	"net/http"

	sysv1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/constants"
	sysclientset "bytetrade.io/web3os/system-server/pkg/generated/clientset/versioned"
	"bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"
	permission "bytetrade.io/web3os/system-server/pkg/permission/v1alpha1"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	MODULE_TAGS       = []string{"message-dispatcher"}
	MODULE_ROUTE      = "/message-dispatcher/v1alpha1"
	PATH_CREATE_EVENT = "/fire-event"
	PATH_LIST_EVENT   = "/list-event"
	PATH_WATCH_EVENT  = "/dispatch-event"

	disableDispatcher = false
)

func init() {
	flag.BoolVar(&disableDispatcher, "disable_dispatcher", false, "disable default message dispatcher")
}

func AddMessageDispatcherToContainer(ctx context.Context,
	sysclientset *sysclientset.Clientset,
	c *restful.Container,
	operator *db.DbOperator,
	ctrlSet *permission.PermissionControlSet,
) error {
	// reconcile crd
	err := reconcileRegistry(ctx, sysclientset, disableDispatcher)
	if err != nil {
		return err
	}

	if !disableDispatcher {
		handler := newHandler(operator, ctrlSet)

		ws := newWebService()

		ws.Route(ws.POST(PATH_CREATE_EVENT).
			To(handler.fireEvent).
			Doc("fire an event").
			Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
			Param(ws.HeaderParameter(api.AccessTokenHeader, "Access token")).
			Returns(http.StatusOK, "Success to fire an event", nil))

		ws.Route(ws.POST(PATH_LIST_EVENT).
			To(handler.listEvents).
			Doc("list events").
			Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
			Param(ws.HeaderParameter(api.AccessTokenHeader, "Access token")).
			Returns(http.StatusOK, "Success to list events", nil))

		ws.Route(ws.POST(PATH_WATCH_EVENT).
			To(handler.dispatchEvent).
			Doc("dispatch an event").
			Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
			Param(ws.HeaderParameter(api.AccessTokenHeader, "Access token")).
			Returns(http.StatusOK, "Success to dispatch an event", nil))

		c.Add(ws)
	}

	return nil
}

func reconcileRegistry(ctx context.Context, sysclientset *sysclientset.Clientset, disable bool) error {
	err := _reconcileRegistry(ctx, sysclientset, disable, DefaultWatcherName, &EVENT_WATCHER)
	if err != nil {
		return err
	}

	err = _reconcileRegistry(ctx, sysclientset, disable, DefaultProviderName, &EVENT_PROVIDER)
	if err != nil {
		return err
	}

	return nil
}

func _reconcileRegistry(ctx context.Context,
	sysclientset *sysclientset.Clientset,
	disable bool,
	name string,
	pr *sysv1alpha1.ProviderRegistry) error {
	oldpr, err := sysclientset.SysV1alpha1().
		ProviderRegistries(constants.MyNamespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		if status, ok := err.(errors.APIStatus); ok || syserrs.As(err, &status) {
			klog.Infof("error is %v, %v", status.Status().Reason, status.Status().Code)
		}

		return err
	}

	if err == nil && oldpr != nil {
		err := sysclientset.SysV1alpha1().
			ProviderRegistries(constants.MyNamespace).
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	if !disable {
		_, err := sysclientset.SysV1alpha1().
			ProviderRegistries(constants.MyNamespace).
			Create(ctx, pr, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func newWebService() *restful.WebService {
	webservice := restful.WebService{}

	webservice.Path(MODULE_ROUTE).
		Produces(restful.MIME_JSON)

	return &webservice
}
