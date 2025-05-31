package migrator

// Dialect is the type of database dialect.
type Dialect string

const (
	DialectClickHouse Dialect = "clickhouse"
	DialectMySQL      Dialect = "mysql"
	DialectPostgres   Dialect = "postgres"
	DialectSQLite3    Dialect = "sqlite3"
	DialectYdB        Dialect = "ydb"
)
