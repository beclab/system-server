package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	apiv1alpha1 "bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/constants"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	NM_URL = "notification-manager-svc.kubesphere-monitoring-system:19093"
)

func init() {
	url := os.Getenv("NM_URL")
	if url != "" {
		NM_URL = url
	}
}

type EventWatcher struct {
	httpClient *resty.Client
}

func NewEventWatcher() *EventWatcher {
	client := resty.New()
	return &EventWatcher{
		httpClient: client.SetTimeout(2 * time.Second),
	}
}

func (e *EventWatcher) DoWatch(event *Event) error {

	pl, err := json.Marshal(event.Data.Payload)
	if err != nil {
		return err
	}

	alert := Alert{
		Labels: KV{
			"namespace": constants.MyNamespace, // MUST to have
			"type":      event.Type,
			"version":   event.Version,
			"payload":   string(pl),
		},
		Annotations: KV{
			"message": event.Data.Message,
		},
	}

	// Set the webhook receiver to notification server for kubeshpere notification manager
	// and send the notification message to Notification Mananger
	// Notification Mananger will pick the webhook receiver to send message
	enableWebhook := true
	notificationServiceUrl := fmt.Sprintf("http://notifications-service.%s/notification/system/push", strings.Replace(constants.MyNamespace, "user-system-", "user-space-", 1))
	request := NotificationManagerRequest{
		Alert: &struct {
			Alerts Alerts `json:"alerts"`
		}{
			Alerts: Alerts{
				alert,
			},
		},
		Receiver: &Receiver{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "notification.kubesphere.io/v2beta2",
				Kind:       "Receiver",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name: "sys-receiver",
				Labels: map[string]string{
					"app":  "notification-manager",
					"type": "global",
				},
			},

			Spec: ReceiverSpec{
				Webhook: &WebhookReceiver{
					Enabled: &enableWebhook,
					URL:     &notificationServiceUrl,
					HTTPConfig: &HTTPClientConfig{
						BasicAuth: &BasicAuth{
							Username: apiv1alpha1.BackendTokenHeader,
							Password: &Credential{
								Value: constants.Nonce,
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	klog.Info("send alert to notification manager, ", string(data))

	postURL := fmt.Sprintf("http://%s/api/v2/notifications", NM_URL)
	res, err := e.httpClient.R().
		SetHeader(restful.HEADER_ContentType, restful.MIME_JSON).
		SetBody(data).
		Post(postURL)

	if err != nil {
		return err
	}

	klog.Info("send alert to notification manager, get response: ", string(res.Body()))

	var nmResp NotificationManagerResponse
	err = json.Unmarshal(res.Body(), &nmResp)
	if err != nil {
		return err
	}

	if nmResp.Status != http.StatusOK {
		return errors.New(nmResp.Message)
	}

	return nil
}
