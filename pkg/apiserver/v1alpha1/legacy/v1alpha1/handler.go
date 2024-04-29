package legacy

import (
	"net/http"
	"net/http/httputil"
	"reflect"

	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	prodiverregistry "bytetrade.io/web3os/system-server/pkg/providerregistry/v1alpha1"
	serviceproxy "bytetrade.io/web3os/system-server/pkg/serviceproxy/v1alpha1"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	"k8s.io/klog/v2"
)

type Handler struct {
	method string
	proxy  *serviceproxy.Proxy
}

func newHandler(method string, registry *prodiverregistry.Registry,
) *Handler {
	proxy := serviceproxy.NewProxy(registry)

	return &Handler{
		method: method,
		proxy:  proxy,
	}
}

func (h *Handler) do(req *restful.Request, resp *restful.Response) {
	klog.Info("proxy ", h.method, " /", req.PathParameter(serviceproxy.ParamSubPath))

	proxyRespIntf, err := h.proxy.ProxyLegacyAPI(req.Request.Context(), h.method, req, resp)
	if err != nil && isNil(proxyRespIntf) {
		klog.Info("proxy error: ", err)
		api.HandleError(resp, req, err)
		return
	}

	if err == nil && isNil(proxyRespIntf) {
		klog.Info("websocket proxy connected")
		return
	}

	switch proxyResp := proxyRespIntf.(type) {
	case *resty.Response:
		dump, err := httputil.DumpRequest(proxyResp.Request.RawRequest, true)
		if err != nil {
			klog.Error("dump request err: ", err)
		} else {
			klog.Info("proxy request: ", string(dump))
		}

		dump, err = httputil.DumpResponse(proxyResp.RawResponse, false)
		if err != nil {
			klog.Error("dump response err: ", err)
		} else {
			klog.Info("proxy response: ", string(dump))
		}

		for h, values := range proxyResp.Header() {
			for _, v := range values {
				resp.Header().Set(h, v)
			}
		}

		for _, c := range proxyResp.Cookies() {
			http.SetCookie(resp, c)
		}

		resp.WriteHeader(proxyResp.StatusCode())
		resp.Write(proxyResp.Body())

	case *serviceproxy.WsProxyResponse:
		resp.WriteHeader(proxyResp.RawResponse.StatusCode)
		resp.Write(proxyResp.Body)
	}

}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
