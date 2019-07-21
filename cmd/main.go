package main

import (
	"fmt"
	"github.com/koluchiy/mysql-study/pkg/connection"
	"github.com/koluchiy/mysql-study/pkg/runner"
	"github.com/koluchiy/mysql-study/pkg/drawer"
	"github.com/koluchiy/mysql-study/pkg/data"
	"github.com/koluchiy/mysql-study/pkg/config"
)

func main() {
	cfg := config.GetConfig()

	loader := data.NewLoader(cfg.DbConfig)

	fmt.Println("Load flow from " + cfg.FlowPath)

	flow, err := loader.Load(cfg.FlowPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	conns := map[string]connection.Connection{}

	for _, c := range flow.Connections {
		conn, err := connection.NewConnection(connection.Config{
			Dsn: c.Dsn,
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		conns[c.Name] = conn
	}

	queries := make([]runner.Query, len(flow.Queries))
	for i, q := range flow.Queries {
		queries[i] = runner.Query{
			Type: q.Type,
			Connection: q.Connection,
			Sql: q.Sql,
			MessageBefore: drawer.NewMessageString(q.MessageBefore),
			MessageAfter: drawer.NewMessageString(q.MessageAfter),
			Async: q.Async,
			Sleep: q.Sleep,
			Timeout: q.Timeout,
		}
	}

	dr := drawer.NewDrawer()

	r := runner.NewRunner(dr)

	err = r.Run(runner.Flow{
		Queries: queries,
		Connections: conns,
	})

	fmt.Println(err)
}
