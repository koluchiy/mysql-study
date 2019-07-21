package config

import (
	"github.com/koluchiy/mysql-study/pkg/data"
	"sync"
	"github.com/namsral/flag"
)

type Config struct {
	FlowPath string
	DbConfig data.Config
}

var cfg Config
var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		var dsn string
		var flowPath string
		flag.StringVar(&dsn, "db_dsn", "", "dsn for database")
		flag.StringVar(&flowPath, "flow_path", "", "path for flow")
		flag.Parse()

		cfg = Config{
			FlowPath: flowPath,
			DbConfig: data.Config{
				DbDsn: dsn,
			},
		}
	})

	return cfg
}
