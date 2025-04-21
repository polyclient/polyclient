module github.com/polyclient/polyclient

go 1.24.2

replace github.com/polyclient/polyclient/drivers/cassandra => ./drivers/cassandra

replace github.com/polyclient/polyclient/drivers/clickhouse => ./drivers/clickhouse

replace github.com/polyclient/polyclient/drivers/cockroachdb => ./drivers/cockroachdb

replace github.com/polyclient/polyclient/drivers/duckdb => ./drivers/duckdb

replace github.com/polyclient/polyclient/drivers/firebird => ./drivers/firebird

replace github.com/polyclient/polyclient/drivers/libsql => ./drivers/libsql

replace github.com/polyclient/polyclient/drivers/mariadb => ./drivers/mariadb

replace github.com/polyclient/polyclient/drivers/mongodb => ./drivers/mongodb

replace github.com/polyclient/polyclient/drivers/mysql => ./drivers/mysql

replace github.com/polyclient/polyclient/drivers/oracle => ./drivers/oracle

replace github.com/polyclient/polyclient/drivers/postgresql => ./drivers/postgresql

replace github.com/polyclient/polyclient/drivers/redis => ./drivers/redis

replace github.com/polyclient/polyclient/drivers/sqlite => ./drivers/sqlite

replace github.com/polyclient/polyclient/drivers/sqlserver => ./drivers/sqlserver

require (
	github.com/bytecodealliance/wasmtime-go/v30 v30.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-playground/validator/v10 v10.25.0
	github.com/polyclient/polyclient/drivers/postgresql v0.0.0-00010101000000-000000000000
	github.com/polyclient/polyclient/drivers/sqlite v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v3 v3.0.0-beta1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.61.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.8.2 // indirect
	modernc.org/sqlite v1.36.2 // indirect
)
