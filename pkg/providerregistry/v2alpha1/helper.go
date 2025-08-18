package v2alpha1

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	headerXForwardedURI = "X-Forwarded-Uri"
)

func ProviderServiceAddr(providerRef string) string {
	token := strings.Split(providerRef, "/")
	if len(token) == 1 {
		return token[0]
	}

	return fmt.Sprintf("%s.%s", token[1], token[0])
}

func ProviderRefFromHost(host string) string {
	token := strings.Split(host, ".")
	if len(token) < 2 {
		return host
	}

	return fmt.Sprintf("%s/%s", token[1], token[0])
}

func ProviderRefName(appName, namespace string) string {
	if len(namespace) == 0 {
		return appName
	}

	return fmt.Sprintf("%s/%s", namespace, appName)
}

// GetXForwardedURI returns the content of the X-Forwarded-URI header, falling back to the start-line request path.
func GetXForwardedURI(req *http.Request) (uri string) {
	uri = req.Header.Get(headerXForwardedURI)

	if len(uri) == 0 {
		return req.URL.String()
	}

	return uri
}
