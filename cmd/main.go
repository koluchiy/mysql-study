package main

import (
	"fmt"
	"github.com/koluchiy/mysql-study/pkg/connection"
	"github.com/koluchiy/mysql-study/pkg/runner"
	"github.com/koluchiy/mysql-study/pkg/drawer"
	"github.com/koluchiy/mysql-study/pkg/data"
)

func main() {
	loader := data.NewLoader()
	flow, err := loader.Load("data/gap-lock.yml")

	if err != nil {
		panic(err)
	}

	conns := map[string]connection.Connection{}

	for _, c := range flow.Connections {
		conn, err := connection.NewConnection(connection.Config{
			Dsn: c.Dsn,
		})

		if err != nil {
			panic(err)
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

	r.Run(runner.Flow{
		Queries: queries,
		Connections: conns,
	})

	//conn1, err := connection.NewConnection(connection.Config{
	//	Dsn: "root:password@tcp(127.0.0.1:3306)/study",
	//})
	//
	//conn2, err := connection.NewConnection(connection.Config{
	//	Dsn: "root:password@tcp(127.0.0.1:3306)/study",
	//})
	//
	//dr := drawer.NewDrawer()
	//
	//r := runner.NewRunner(dr)
	//err = r.Run(runner.Flow{
	//	Connections: map[string]connection.Connection{
	//		"conn1": conn1,
	//		"conn2": conn2,
	//	},
	//	Queries: []runner.Query{
	//		{
	//			Type: runner.QueryTypeQuery,
	//			Connection: "conn1",
	//			Sql: "select * from users",
	//		},
	//		{
	//			Type: runner.QueryTypeExec,
	//			Connection: "conn2",
	//			Sql: "begin",
	//		},
	//		{
	//			Type: runner.QueryTypeExec,
	//			Connection: "conn2",
	//			Sql: "insert into users(title) values('user3')",
	//		},
	//		{
	//			Type: runner.QueryTypeQuery,
	//			Connection: "conn1",
	//			Sql: "select * from users",
	//			MessageAfter: drawer.NewMessageString("We still not see new row"),
	//		},
	//		{
	//			Type: runner.QueryTypeExec,
	//			Connection: "conn2",
	//			Sql: "commit",
	//		},
	//		{
	//			Type: runner.QueryTypeQuery,
	//			Connection: "conn1",
	//			Sql: "select * from users",
	//			MessageAfter: drawer.NewMessageString("transaction commited, we see new row"),
	//		},
	//	},
	//})

	fmt.Println(err)

	//res, err := conn.QueryContext(context.Background(), "select * from users")
	//
	//fmt.Println(err, res)
}
