package constants

import (
	"os"

	"github.com/google/uuid"
)

const (
	APIServerListenAddress    = ":80"
	KubeSphereClientAttribute = "ksclient"
	AuthorizationTokenKey     = "X-Authorization"
	BflUserKey                = "X-BFL-USER"
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
