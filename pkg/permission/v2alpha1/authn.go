package v2alpha1

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bytetrade.io/web3os/system-server/pkg/constants"
	"github.com/brancz/kube-rbac-proxy/pkg/authn"
	"github.com/golang-jwt/jwt"
	"github.com/jellydator/ttlcache/v3"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var _ authenticator.Request = (*lldapTokenAuthenticator)(nil)

type Claims struct {
	jwt.StandardClaims
	// Private Claim Names
	// Username user identity, deprecated field
	Username string `json:"username,omitempty"`

	Groups []string `json:"groups,omitempty"`
	Mfa    int64    `json:"mfa,omitempty"`
}

type lldapTokenAuthenticator struct {
	tokenCache *ttlcache.Cache[string, *Claims]
}

// AuthenticateRequest implements authenticator.Request.
func (l *lldapTokenAuthenticator) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
	token := req.Header.Get(constants.AuthorizationTokenKey)
	if token == "" {
		cookie, err := req.Cookie(constants.AuthTokenCookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				return nil, false, nil // No token found
			}
			return nil, false, fmt.Errorf("error retrieving cookie: %w", err)
		}

		token = cookie.Value
	}

	claims := l.tokenCache.Get(token)
	if claims != nil {
		return &authenticator.Response{
			User: &user.DefaultInfo{
				Name:   claims.Value().Username,
				Groups: claims.Value().Groups,
				UID:    claims.Value().Username,
			},
		}, true, nil
	}

	// TODO:
	return nil, false, nil
}

func UnionAllAuthenticators(ctx context.Context, cfg *authn.AuthnConfig, kubeClient kubernetes.Interface) (authenticator.Request, error) {
	var authenticator authenticator.Request

	// If OIDC configuration provided, use oidc authenticator
	if cfg.OIDC.IssuerURL != "" {
		oidcAuthenticator, err := authn.NewOIDCAuthenticator(ctx, cfg.OIDC)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate OIDC authenticator: %w", err)
		}

		go oidcAuthenticator.Run(ctx)
		authenticator = oidcAuthenticator
	} else {
		//Use Delegating authenticator
		klog.Infof("Valid token audiences: %s", strings.Join(cfg.Token.Audiences, ", "))

		tokenClient := kubeClient.AuthenticationV1()
		delegatingAuthenticator, err := authn.NewDelegatingAuthenticator(tokenClient, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate delegating authenticator: %w", err)
		}

		go delegatingAuthenticator.Run(ctx)
		authenticator = delegatingAuthenticator
	}

	return union.New(&lldapTokenAuthenticator{ttlcache.New(
		ttlcache.WithTTL[string, *Claims](time.Minute*5),
		ttlcache.WithCapacity[string, *Claims](1000),
	)}, authenticator), nil
}
