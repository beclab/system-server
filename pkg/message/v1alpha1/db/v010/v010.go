package v010

import "bytetrade.io/web3os/system-server/pkg/message/v1alpha1/db"

func init() {
	db.AddToDbInit([]db.InitSript{
		create_table,
	})
}

var (
	create_table = db.InitSript{
		Info: "Initialize message event tables",
		SQL: `
create table if not exists user_events(
	id integer primary key autoincrement,
	event_type varchar(50) not null,
	raw_message text not null,
	create_time datetime default CURRENT_TIMESTAMP
);

create index if not exists event_time on user_events(create_time);
`}
)
