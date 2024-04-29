package message

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"
	serviceproxy "bytetrade.io/web3os/system-server/pkg/serviceproxy/v1alpha1"

	"k8s.io/klog/v2"
)

type EventProvider struct {
	dbOperator *db.DbOperator
}

func NewEventProvider(operator *db.DbOperator) *EventProvider {
	db.InitDB(operator)

	return &EventProvider{
		dbOperator: operator,
	}
}

func (p *EventProvider) createEvent(ctx context.Context, event *Event) error {

	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}

	result, err := p.dbOperator.DB.NamedExecContext(ctx,
		"insert into user_events(event_type, raw_message)values(:event_type, :raw_message)",
		&db.UserEvents{
			EventType:  event.Type,
			RawMessage: string(raw),
		},
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	klog.Info("insert event with id, ", id)
	return nil
}

func (p *EventProvider) listEvent(ctx context.Context, param *serviceproxy.ListOpParam) ([]db.UserEvents, error) {
	sql := "select * from user_events"
	if len(param.Filters) > 0 {
		sql += " where "
	}

	filterStrs := []string{}
	args := []interface{}{}
	index := 1
	for k, f := range param.Filters {
		switch k {
		case "create_time":
			if len(f) == 1 {
				filterStrs = append(filterStrs, fmt.Sprintf("create_time <= $%d", index))
				t, e := stringToTime(f[0])
				if e != nil {
					return nil, e
				}
				args = append(args, t)
				index++
			} else {
				filterStrs = append(filterStrs, fmt.Sprintf("create_time >= $%d and create_time <= $%d", index, index+1))
				t, e := stringToTime(f[0])
				if e != nil {
					return nil, e
				}
				args = append(args, t)

				t, e = stringToTime(f[1])
				if e != nil {
					return nil, e
				}
				args = append(args, t)
				index += 2
			}
		default:
			inArgs := []string{}
			for _, subf := range f {
				args = append(args, subf)
				inArgs = append(inArgs, fmt.Sprintf("$%d", index))
				index++
			}

			filterStrs = append(filterStrs, fmt.Sprintf("%s in (%s)", k, strings.Join(inArgs, " , ")))
		}
	}

	sql += strings.Join(filterStrs, " and ")
	sql += " order by create_time desc "
	sql += fmt.Sprintf(" limit %d, %d", param.Page.Offset, param.Page.Limit)

	klog.Info("list sql, ", sql, " args: ", args)

	list := make([]db.UserEvents, 0)
	err := p.dbOperator.DB.SelectContext(ctx, &list, sql, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func stringToTime(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}
