# migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" drop -f
migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" up
# migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" down

# go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest add variable to PATH
# sqlc generate

XGROUP DESTROY events worker-group
XGROUP CREATE events worker-group 0 MKSTREAM
XGROUP CREATE events worker-group $ MKSTREAM

# PostgreSQL
# migrate -path ./migrations -database "postgres://user:password@host:port/dbname?sslmode=disable" up

# SQLite
# migrate -path ./migrations -database "sqlite3://absolute/path/to/file.db" up

# Microsoft SQL Server
# migrate -path ./migrations -database "sqlserver://user:password@host:port?database=dbname" up

# CockroachDB (Postgres-compatible)
# migrate -path ./migrations -database "postgres://user:password@host:port/dbname?sslmode=disable" up
# migrate -path ./migrations -database "cockroachdb://user:password@host:port/dbname?sslmode=disable" up

# ClickHouse
# migrate -path ./migrations -database "clickhouse://default:123456@localhost:9000/test" up


# Cassandra
# migrate -path ./migrations -database "cassandra://user:password@host:port/keyspace" up
# migrate -path ./migrations -database "cassandra://127.0.0.1:9042/testkeyspace?username=cassandra&password=cassandra" up

# ScyllaDB
# migrate -path ./migrations -database "scylla://user:password@host:port/keyspace" up
# migrate -path ./migrations -database "cassandra://127.0.0.1:9042/mykeyspace?username=scylla&password=scylla" up

# go install -tags 'postgres mysql cassandra clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# https://chatgpt.com/c/68aed56d-61b0-8323-a0d4-ecc0ad99c5ce
# add variable to PATH