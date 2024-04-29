package message

import (
	"encoding/json"
	"errors"

	sysv1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api"
	"bytetrade.io/web3os/system-server/pkg/apiserver/v1alpha1/api/response"
	"bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"
	permission "bytetrade.io/web3os/system-server/pkg/permission/v1alpha1"
	serviceproxy "bytetrade.io/web3os/system-server/pkg/serviceproxy/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/utils"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
)

type Handler struct {
	provider       *EventProvider
	watcher        *EventWatcher
	permissionCtrl *permission.PermissionControlSet
}

func newHandler(operator *db.DbOperator,
	ctrlSet *permission.PermissionControlSet,
) *Handler {

	provider := NewEventProvider(operator)
	watcher := NewEventWatcher()
	return &Handler{
		provider:       provider,
		watcher:        watcher,
		permissionCtrl: ctrlSet,
	}
}

func (h *Handler) fireEvent(req *restful.Request, resp *restful.Response) {
	var eventReq serviceproxy.ProxyRequest
	event, err := h.getEventFromRequst(req, &eventReq)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	klog.Infof("fire event: %s", utils.PrettyJSON(event))

	err = h.provider.createEvent(req.Request.Context(), event)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	response.SuccessNoData(resp)
}

func (h *Handler) dispatchEvent(req *restful.Request, resp *restful.Response) {
	var eventReq serviceproxy.DispatchRequest
	event, err := h.getEventFromRequst(req, &eventReq)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	klog.Infof("dispatch event: %s", utils.PrettyJSON(event))

	err = h.watcher.DoWatch(event)
	if err != nil {
		// internal api, tell the invoker if the api error
		api.HandleError(resp, req, err)
		return
	}

	response.SuccessNoData(resp)
}

func (h *Handler) listEvents(req *restful.Request, resp *restful.Response) {
	var eventReq serviceproxy.ProxyRequest
	err := req.ReadEntity(&eventReq)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	jsonData, err := json.Marshal(eventReq.Param)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	var param serviceproxy.ListOpParam
	err = json.Unmarshal(jsonData, &param)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	events, err := h.provider.listEvent(req.Request.Context(), &param)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	retEvents := make([]Event, 0)
	for _, e := range events {
		var ev Event
		err := json.Unmarshal([]byte(e.RawMessage), &ev)
		if err != nil {
			response.HandleError(resp, err)
			return
		}

		retEvents = append(retEvents, ev)
	}

	response.Success(resp, retEvents)
}

func (h *Handler) getEventFromRequst(req *restful.Request, eventReqPtr interface{}) (*Event, error) {
	err := req.ReadEntity(eventReqPtr)
	if err != nil {
		return nil, err
	}

	klog.Infof("event reqeust, %s", utils.PrettyJSON(eventReqPtr))
	switch eventReq := eventReqPtr.(type) {
	case *serviceproxy.DispatchRequest:
		// TODO: check result
		_, err := permission.ValidateAccessToken(eventReq.Token,
			eventReq.Op,
			eventReq.DataType,
			eventReq.Version,
			eventReq.Group,
			h.permissionCtrl)
		if err != nil {
			return nil, err
		}

		if !utils.ListContains(sysv1alpha1.WatcherSupportedOPs, eventReq.Op) {
			return nil, errors.New("unsupported: method not provided")
		}

		if eventReq.DataType != sysv1alpha1.Event {
			return nil, errors.New("unsupported: wrong data type")
		}

		if eventReq.Group != GroupID {
			return nil, errors.New("unsupported: group error")
		}

		event, err := convertDataToEvent(eventReq.Data)
		if err != nil {
			return nil, err
		}

		return event, nil
	case *serviceproxy.ProxyRequest:
		_, err := permission.ValidateAccessToken(eventReq.Token,
			eventReq.Op,
			eventReq.DataType,
			eventReq.Version,
			eventReq.Group,
			h.permissionCtrl)
		if err != nil {
			return nil, err
		}

		if eventReq.DataType != sysv1alpha1.Event {
			return nil, errors.New("unsupported: wrong data type")
		}

		if eventReq.Group != GroupID {
			return nil, errors.New("unsupported: group error")
		}

		event, err := convertDataToEvent(eventReq.Data)
		if err != nil {
			return nil, err
		}

		return event, nil
	}

	return nil, errors.New("unknown request type")
}

func convertDataToEvent(data interface{}) (*Event, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var event Event
	err = json.Unmarshal(jsonData, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
