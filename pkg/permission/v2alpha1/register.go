package v2alpha1

import (
	"net/http"

	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api/response"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	MODULE_TAGS  = []string{"permission-control"}
	MODULE_ROUTE = "/permission/v2alpha1"
)

func AddPermissionControlToContainer(
	c *restful.Container,
	authenticator authenticator.Request,
	kubeconfig *rest.Config,
) error {
	client := kubernetes.NewForConfigOrDie(kubeconfig)
	handler := &handler{kubeClient: client}
	ws := newWebService()

	requireAuth := func(f restful.RouteFunction) restful.RouteFunction {
		return func(req *restful.Request, resp *restful.Response) {
			handlerFunc := func(rw http.ResponseWriter, r *http.Request) {
				f(req, resp)
			}

			handlerFunc = WithUserHeader(handlerFunc)
			handlerFunc = WithAuthentication(authenticator, nil, handlerFunc)

			handlerFunc(resp, req.Request)
		}
	}

	ws.Route(ws.POST("/register").
		To(requireAuth(handler.register)).
		Doc("register an app provider binding").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Returns(http.StatusOK, "Success to register a invoker", &RegisterResp{}))

	ws.Route(ws.POST("/unregister").
		To(requireAuth(handler.unregister)).
		Doc("unregister an app provider binding").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Returns(http.StatusOK, "Success to unregister a invoker", &response.Response{}))

	c.Add(ws)

	return nil
}

func newWebService() *restful.WebService {
	webservice := restful.WebService{}

	webservice.Path(MODULE_ROUTE).
		Produces(restful.MIME_JSON)

	return &webservice
}
