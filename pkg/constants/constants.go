package constants

import (
	"os"

	"github.com/google/uuid"
)

const (
	ProxyServerListenAddress  = ":28080"
	APIServerListenAddress    = ":84"
	KubeSphereClientAttribute = "ksclient"
	AuthorizationTokenKey     = "X-Authorization"
	BflUserKey                = "X-BFL-USER"
	AuthTokenCookieName       = "auth_token"
)

var (
	MyNamespace string
	Owner       string
)

var (
	Nonce = uuid.New().String()
)

func init() {
	MyNamespace = os.Getenv("MY_NAMESPACE")
	Owner = os.Getenv("OWNER")
}
