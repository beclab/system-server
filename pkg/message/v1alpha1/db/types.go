package db

import (
	"time"
)

type UserEvents struct {
	Id         int
	EventType  string    `db:"event_type"`
	RawMessage string    `db:"raw_message"`
	CreateTime time.Time `db:"create_time"`
}
