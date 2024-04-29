package constants

import (
	"os"

	"github.com/google/uuid"
)

const (
	APIServerListenAddress    = ":80"
	KubeSphereClientAttribute = "ksclient"
	AuthorizationTokenKey     = "X-Authorization"
)

var (
	MyNamespace string
)

var (
	Nonce = uuid.New().String()
)

func init() {
	MyNamespace = os.Getenv("MY_NAMESPACE")
}
