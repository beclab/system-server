package permission

import (
	"errors"
	"net/http"

	sysv1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api/response"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/client-go/rest"
)

var (
	MODULE_TAGS  = []string{"permission-control"}
	MODULE_ROUTE = "/permission/v1alpha1"
)

func AddPermissionControlToContainer(c *restful.Container,
	ctrlSet *PermissionControlSet,
	kubeconfig *rest.Config,
) error {
	handler := newHandler(ctrlSet, kubeconfig)

	ws := newWebService()
	ws.Route(ws.POST("/access").
		To(handler.auth).
		Doc("request a data access token").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Returns(http.StatusOK, "Success to get a token", &AccessTokenResponse{}))

	ws.Route(ws.POST("/register").
		To(handler.register).
		Doc("register a data invoker").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Param(ws.HeaderParameter(api.AuthorizationTokenHeader, "Auth token")).
		Returns(http.StatusOK, "Success to register a invoker", &RegisterResp{}))

	ws.Route(ws.POST("/unregister").
		To(handler.unregister).
		Doc("unregister a data invoker").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Param(ws.HeaderParameter(api.AuthorizationTokenHeader, "Auth token")).
		Returns(http.StatusOK, "Success to unregister a invoker", &response.Response{}))

	ws.Route(ws.GET("/nonce").
		To(handler.nonce).
		Doc("get backend request call nonce").
		Metadata(restfulspec.KeyOpenAPITags, MODULE_TAGS).
		Returns(http.StatusOK, "Success to get nonce", ""))

	c.Add(ws)

	return nil
}

func newWebService() *restful.WebService {
	webservice := restful.WebService{}

	webservice.Path(MODULE_ROUTE).
		Produces(restful.MIME_JSON)

	return &webservice
}

func ValidateAccessTokenWithRequest(token string, op string, req *restful.Request, ctrlSet *PermissionControlSet) (string, error) {
	datatype := req.PathParameter(api.ParamDataType)
	version := req.PathParameter(api.ParamVersion)
	group := req.PathParameter(api.ParamGroup)

	return ValidateAccessToken(token, op, datatype, version, group, ctrlSet)
}

func ValidateAccessToken(token string, op, datatype, version, group string, ctrlSet *PermissionControlSet) (string, error) {
	permReq, err := ctrlSet.Mgr.getPermWithToken(token)
	if err != nil {
		return "", err
	}

	accReq := sysv1alpha1.PermissionRequire{
		Group:    group,
		DataType: datatype,
		Version:  version,
		Ops: []string{
			op,
		},
	}

	if permReq.Include(&accReq, true) {
		return permReq.AppKey, nil
	}

	return "", errors.New("data access denied")
}
