package message

import (
	"context"
	"flag"
	"testing"

	"bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"
	serviceproxy "bytetrade.io/web3os/system-server/pkg/serviceproxy/v1alpha1"
	"bytetrade.io/web3os/system-server/pkg/utils"
)

func newTestOperator() *EventProvider {
	flag.Set("db", "/tmp/test.db")
	opr := db.NewDbOperator()
	return NewEventProvider(opr)
}

func TestCreateEvent(t *testing.T) {
	provider := newTestOperator()
	err := provider.createEvent(context.TODO(), &Event{
		Type:    "notify-test",
		Version: "v1",
		Data: EventData{
			Message: "test event",
			Payload: "test payload",
		},
	})

	if err != nil {
		t.Fatal("create event error, ", err)
	}
}

func TestListEvent(t *testing.T) {
	provider := newTestOperator()

	evs, err := provider.listEvent(context.TODO(), &serviceproxy.ListOpParam{
		Filters: map[string][]string{
			"event_type": {
				"hahah",
				"notify-test",
			},
			"create_time": {
				"2022-12-30 11:00:00",
				"2022-12-31 01:00:00",
			},
		},
		Page: serviceproxy.Pagination{
			Offset: 0,
			Limit:  100,
		},
	})

	if err != nil {
		t.Fatal("list event error, ", err)
	}

	t.Log("result len: ", len(evs), " items: ", utils.PrettyJSON(evs))
}
