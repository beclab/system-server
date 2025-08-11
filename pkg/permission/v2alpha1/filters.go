package v2alpha1

import (
	"fmt"
	"net/http"
	"net/url"

	providerv2alpha1 "bytetrade.io/web3os/system-server/pkg/providerregistry/v2alpha1"
	"github.com/brancz/kube-rbac-proxy/pkg/authz"
	"github.com/brancz/kube-rbac-proxy/pkg/proxy"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

func WithAuthentication(
	authReq authenticator.Request,
	audiences []string,
	handler http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if len(audiences) > 0 {
			ctx = authenticator.WithAudiences(ctx, audiences)
			req = req.WithContext(ctx)
		}

		res, ok, err := authReq.AuthenticateRequest(req)
		if err != nil {
			klog.Errorf("Unable to authenticate the request due to an error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req = req.WithContext(request.WithUser(req.Context(), res.User))
		handler.ServeHTTP(w, req)
	}
}

func WithAuthorization(
	authz authorizer.Authorizer,
	cfg *authz.Config,
	handler http.HandlerFunc,
) http.HandlerFunc {
	getRequestAttributes := func(u user.Info, r *http.Request) []authorizer.Attributes {
		allAttrs := proxy.
			NewKubeRBACProxyAuthorizerAttributesGetter(cfg).
			GetRequestAttributes(u, r)

		for i, attrs := range allAttrs {
			if attrs.GetPath() != "" && !attrs.IsResourceRequest() {
				// for non-resource requests, setup the provider reference
				uri := providerv2alpha1.GetXForwardedURI(r)
				requestUrl, err := url.Parse(uri)
				if err != nil {
					klog.Errorf("failed to parse X-Forwarded-URI: %v", err)
					return nil
				}
				hostStr := requestUrl.Host
				ref := providerv2alpha1.ProviderRefFromHost(hostStr)

				a := authorizer.AttributesRecord{
					User:            attrs.GetUser(),
					Verb:            attrs.GetVerb(),
					Namespace:       attrs.GetNamespace(),
					APIGroup:        attrs.GetAPIGroup(),
					APIVersion:      attrs.GetAPIVersion(),
					Resource:        ref,
					Subresource:     attrs.GetSubresource(),
					Name:            attrs.GetName(),
					ResourceRequest: attrs.IsResourceRequest(),
					Path:            attrs.GetPath(),
				}
				allAttrs[i] = a
			}
		}

		return allAttrs
	}

	return func(w http.ResponseWriter, req *http.Request) {
		u, ok := request.UserFrom(req.Context())
		if !ok {
			http.Error(w, "user not in context", http.StatusBadRequest)
			return
		}

		// Get authorization attributes
		allAttrs := getRequestAttributes(u, req)
		if len(allAttrs) == 0 {
			msg := "Bad Request. The request or configuration is malformed."
			klog.V(2).Info(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		for _, attrs := range allAttrs {
			// Authorize
			authorized, reason, err := authz.Authorize(req.Context(), attrs)
			if err != nil {
				msg := fmt.Sprintf("Authorization error (user=%s, verb=%s, resource=%s, subresource=%s)", u.GetName(), attrs.GetVerb(), attrs.GetResource(), attrs.GetSubresource())
				klog.Errorf("%s: %s", msg, err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			if authorized != authorizer.DecisionAllow {
				msg := fmt.Sprintf("Forbidden (user=%s, verb=%s, resource=%s, subresource=%s)", u.GetName(), attrs.GetVerb(), attrs.GetResource(), attrs.GetSubresource())
				klog.V(2).Infof("%s. Reason: %q.", msg, reason)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
		}

		handler.ServeHTTP(w, req)
	}
}
