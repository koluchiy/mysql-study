Play flow of mysql queries in different connections

Demonstate myasl features like transactions, lock etc.

# Usage
Export connection to your database

`export DB_DSN="root:password@tcp(127.0.0.1:3306)/study"`

Run with path to your flow yaml file

`go run cmd/main.go -flow_path=data/flow.yml`