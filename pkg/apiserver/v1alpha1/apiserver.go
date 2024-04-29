package apiserver

import (
	"context"
	"net/http"
	"time"

	legacy "bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/legacy/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/constants"
	sysclientset "bytetrade.io/web3os/system-server/pkg/generated/clientset/versioned"
	message "bytetrade.io/web3os/system-server/pkg/message/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"
	permission "bytetrade.io/web3os/system-server/pkg/permission/v1alpha1"
	prodiverregistry "bytetrade.io/web3os/system-server/pkg/providerregistry/v1alpha1"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

// APIServer represents an API server for system.
type APIServer struct {
	Server *http.Server

	// RESTful Server
	container *restful.Container

	serverCtx  context.Context
	dbOperator *db.DbOperator
}

// New constructs a new APIServer.
func New(ctx context.Context) (*APIServer, error) {
	server := &http.Server{
		Addr: constants.APIServerListenAddress,
	}

	operator := db.NewDbOperator()

	return &APIServer{
		Server:     server,
		container:  restful.NewContainer(),
		serverCtx:  ctx,
		dbOperator: operator,
	}, nil
}

// PrepareRun do prepares for API server.
func (s *APIServer) PrepareRun(kubeconfig *rest.Config, sysclientset *sysclientset.Clientset) error {
	s.container.Filter(logRequestAndResponse)
	s.container.Router(restful.CurlyRouter{})
	s.container.RecoverHandler(func(panicReason interface{}, httpWriter http.ResponseWriter) {
		logStackOnRecover(panicReason, httpWriter)
	})

	registry := prodiverregistry.NewRegistry(sysclientset)
	ctrlSet := permission.PermissionControlSet{
		Ctrl: permission.NewPermissionControl(sysclientset),
		Mgr:  permission.NewAccessManager(),
	}

	// use the server context for goroutine in background
	utilruntime.Must(addServiceToContainer(s.serverCtx, s.container, kubeconfig, registry, &ctrlSet))
	utilruntime.Must(message.AddMessageDispatcherToContainer(s.serverCtx, sysclientset, s.container, s.dbOperator, &ctrlSet))
	utilruntime.Must(permission.AddPermissionControlToContainer(s.container, &ctrlSet, kubeconfig))
	utilruntime.Must(legacy.AddLegacyAPIToContainer(s.container, registry))

	s.installAPIDocs()

	s.Server.Handler = s.container

	return nil
}

// Run running a server.
func (s *APIServer) Run() error {
	shutdownCtx, cancel := context.WithTimeout(s.serverCtx, 2*time.Minute)
	defer func() {
		cancel()
		s.dbOperator.Close()
	}()

	go func() {
		<-s.serverCtx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
		klog.Info("shutdown apiserver for system-server")
	}()

	klog.Info("starting apiserver for system-server,", "listen on ", constants.APIServerListenAddress)
	return s.Server.ListenAndServe()
}

func (s *APIServer) installAPIDocs() {
	config := restfulspec.Config{
		WebServices:                   s.container.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/system-server/v1alpha1/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	s.container.Add(restfulspec.NewOpenAPIService(config))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "system-server",
			Description: "system server, service bus for system",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "bytetrade",
					Email: "dev@bytetrade.io",
					URL:   "http://bytetrade.io",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache License 2.0",
					URL:  "http://www.apache.org/licenses/LICENSE-2.0",
				},
			},
			Version: "0.1.0",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "system-server",
		Description: "Web 3 OS system-server"}}}
}
