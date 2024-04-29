package message

import (
	"os"

	sysv1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/constants"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupID             = "message-disptahcer.system-server"
	EventVersion        = "v1"
	DefaultProviderName = "default-event-provider"
	DefaultWatcherName  = "default-event-watcher"
)

func endpoint() string {
	ep := os.Getenv("ENDPOINT")
	if ep != "" {
		return ep
	}

	return "localhost" // "system-server." + constants.MyNamespace
}

type Event struct {
	Type    string    `json:"type"`
	Version string    `json:"version"`
	Data    EventData `json:"data"`
}

type EventData struct {
	Message string      `json:"msg"`
	Payload interface{} `json:"payload"`
}

var (
	EVENT_PROVIDER = sysv1alpha1.ProviderRegistry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "sys.bytetrade.io/v1alpha1",
			Kind:       "ProviderRegistry",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      DefaultProviderName,
			Namespace: constants.MyNamespace,
		},
		Spec: sysv1alpha1.ProviderRegistrySpec{
			Description: "default event provider in system-server",
			Group:       GroupID,
			Kind:        sysv1alpha1.Provider,
			DataType:    sysv1alpha1.Event,
			Version:     EventVersion,
			Endpoint:    endpoint(),
			OpApis: []sysv1alpha1.OpApisItem{
				{
					Name: sysv1alpha1.Create,
					URI:  MODULE_ROUTE + PATH_CREATE_EVENT,
				},
				{
					Name: sysv1alpha1.List,
					URI:  MODULE_ROUTE + PATH_LIST_EVENT,
				},
			},
		},
		Status: sysv1alpha1.ProviderRegistryStatus{
			State: sysv1alpha1.Active,
		},
	}

	EVENT_WATCHER = sysv1alpha1.ProviderRegistry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "sys.bytetrade.io/v1alpha1",
			Kind:       "ProviderRegistry",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      DefaultWatcherName,
			Namespace: constants.MyNamespace,
		},
		Spec: sysv1alpha1.ProviderRegistrySpec{
			Description: "default event provider in system-server",
			Group:       GroupID,
			Kind:        sysv1alpha1.Watcher,
			DataType:    sysv1alpha1.Event,
			Version:     EventVersion,
			Endpoint:    endpoint(),
			Callbacks: []sysv1alpha1.Callback{
				{
					Op:  sysv1alpha1.Create,
					URI: MODULE_ROUTE + PATH_WATCH_EVENT,
					Filters: map[string][]string{
						"type": {
							"notification",
						},
					},
				},
			},
		},
		Status: sysv1alpha1.ProviderRegistryStatus{
			State: sysv1alpha1.Active,
		},
	}
)
