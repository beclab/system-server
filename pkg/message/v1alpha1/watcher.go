package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	apiv1alpha1 "bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/constants"

	"github.com/go-resty/resty/v2"
)

//var (
//	NM_URL = "notification-manager-svc.kubesphere-monitoring-system:19093"
//)
//
//func init() {
//	url := os.Getenv("NM_URL")
//	if url != "" {
//		NM_URL = url
//	}
//}

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

	notificationServiceUrl := fmt.Sprintf("http://notifications-server.%s/notification/system/push", strings.Replace(constants.MyNamespace, "user-system-", "user-space-", 1))

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(alert)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, notificationServiceUrl, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apiv1alpha1.BackendTokenHeader, constants.Nonce)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	_, err = doHttpRequest(context.TODO(), client, req)

	if err != nil {
		return err
	}

	return nil
}

func doHttpRequest(ctx context.Context, client *http.Client, request *http.Request) ([]byte, error) {

	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := ""
		if len(body) > 0 {
			msg = string(body)
		}
		return body, fmt.Errorf("%d, %s", resp.StatusCode, msg)
	}

	return body, nil
}
