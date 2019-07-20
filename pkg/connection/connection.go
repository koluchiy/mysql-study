package connection

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Dsn string
}

type Connection interface {
	ExecContext(ctx context.Context, query string) error
	QueryContext(ctx context.Context, query string) (QueryResult, error)
}

func NewConnection(cfg Config) (*connection, error) {
	db, err := sql.Open("mysql", cfg.Dsn)

	if err != nil {
		return nil ,err
	}

	conn, err := db.Conn(context.Background())

	if err != nil {
		return nil ,err
	}

	var connectionID string
	connStmt := `SELECT CONNECTION_ID()`
	err = conn.QueryRowContext(context.Background(), connStmt).Scan(&connectionID)
	if err != nil {
		return nil, err
	}

	instance := &connection{
		db: db,
		config: cfg,
		conn: conn,
		connectionID: connectionID,
	}

	return instance, nil
}

type connection struct {
	config Config
	db *sql.DB
	conn *sql.Conn
	connectionID string
}

type QueryResult struct {
	Columns []string
	Rows [][]sql.RawBytes
}

func (c *connection) ExecContext(ctx context.Context, query string) error {
	_, err := c.conn.ExecContext(ctx, query)

	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			c.db.Exec("KILL QUERY ?", c.connectionID)
			c.conn, _ = c.db.Conn(context.Background())
		}

		return err
	}

	return err
}

func (c *connection) QueryContext(ctx context.Context, query string) (QueryResult, error) {
	rows, err := c.conn.QueryContext(ctx, query)

	if err != nil {
		return QueryResult{}, err
	}

	columns, err := rows.Columns()

	if err != nil {
		return QueryResult{}, err
	}

	result := QueryResult{
		Columns: columns,
	}

	for rows.Next() {
		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return QueryResult{}, err
		}

		result.Rows = append(result.Rows, values)
	}

	return result, nil
}