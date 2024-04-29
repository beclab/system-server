package db

import "k8s.io/klog/v2"

var (
	init_db_script = []InitSript{}
)

type InitSript struct {
	Info string
	SQL  string
}

func AddToDbInit(sqls []InitSript) {
	init_db_script = append(init_db_script, sqls...)
}

func InitDB(operator *DbOperator) {
	for _, s := range init_db_script {
		klog.Info("Init or update db, ", s.Info)
		operator.DB.MustExec(s.SQL)
	}
}
