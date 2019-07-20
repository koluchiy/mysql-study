package runner

import (
	"github.com/koluchiy/mysql-study/pkg/connection"
	"github.com/pkg/errors"
	"context"
	"github.com/koluchiy/mysql-study/pkg/drawer"
	"time"
)

type Runner interface {
	Run(flow Flow) error
}

type runner struct {
	drawer drawer.Drawer
}

func NewRunner(drawer drawer.Drawer) *runner {
	instance := &runner{
		drawer: drawer,
	}

	return instance
}

const QueryTypeExec = "exec"
const QueryTypeQuery = "query"

type Query struct {
	Type string
	Sql string
	Connection string
	MessageBefore drawer.Message
	MessageAfter drawer.Message
	Async bool
	Sleep int
	Timeout int
}

type Flow struct {
	Connections map[string]connection.Connection
	Queries []Query
}

func wrapMessageConnection(m drawer.Message, conn string) drawer.Message {
	m.Msg = conn + ": " + m.Msg

	return m
}

func (r *runner) Run(flow Flow) error {
	ch := make(chan Query)
	errCh := make(chan error, 10)
	drawCh := make(chan interface{})
	doneCh := make(chan bool)

	go r.draw(drawCh)
	go r.listen(flow, ch, errCh, drawCh, doneCh)

	for _, q := range flow.Queries {
		if q.Sleep > 0 {
			time.Sleep(time.Duration(q.Sleep) * time.Second)
		}
		if !q.MessageBefore.IsEmpty() {
			r.drawer.DrawMessage(wrapMessageConnection(q.MessageBefore, q.Connection))
		}
		_, ok := flow.Connections[q.Connection]
		if !ok {
			return errors.New("connection not found")
		}

		time.Sleep(time.Second)
		ch <- q
		<-doneCh
	}

	time.Sleep(20 * time.Second)

	return nil
}

func (r *runner) draw(ch chan interface{}) {
	for d := range ch {
		qr, ok := d.(connection.QueryResult)
		if ok {
			r.drawer.DrawQueryResult(qr)
		}

		m, ok := d.(drawer.Message)
		if ok {
			r.drawer.DrawMessage(m)
		}
	}
}

func (r *runner) runQuery(conn connection.Connection, q Query, errCh chan error, drawCh chan interface{}) {
	if q.Type == QueryTypeQuery {
		result, err := conn.QueryContext(context.Background(), q.Sql)
		if err != nil {
			drawCh <- wrapMessageConnection(drawer.NewMessageString("err: " + err.Error()), q.Connection)
			errCh <- err
		} else {
			drawCh <- result
		}
	} else if q.Type == QueryTypeExec {
		var ctx context.Context
		var cancel context.CancelFunc
		if q.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(q.Timeout) * time.Second)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		err := conn.ExecContext(ctx, q.Sql)
		if err != nil {
			drawCh <- wrapMessageConnection(drawer.NewMessageString("err:" + err.Error()), q.Connection)
			errCh <- err
		} else {
			drawCh <- drawer.NewMessageString(q.Connection + ": success")
		}
	} else {
		errCh <- errors.New("undefined query type")
	}

	if !q.MessageAfter.IsEmpty() {
		drawCh <- wrapMessageConnection(q.MessageAfter, q.Connection)
	}
}

func (r *runner) listen(flow Flow, ch chan Query, errCh chan error, drawCh chan interface{}, doneCh chan bool) {
	for q := range ch {
		conn := flow.Connections[q.Connection]

		drawCh <- drawer.NewMessageString(q.Connection + ": " + q.Sql)

		if q.Async {
			go r.runQuery(conn, q, errCh, drawCh)
			doneCh <- true
		} else {
			r.runQuery(conn, q, errCh, drawCh)
			doneCh <- true
		}
	}
}
