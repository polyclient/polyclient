#!/usr/bin/env bash

rm -f go.work go.work.sum

go work init
go work use .
go work use ./drivers/cassandra
go work use ./drivers/clickhouse
go work use ./drivers/cockroachdb
go work use ./drivers/duckdb
go work use ./drivers/firebird
go work use ./drivers/libsql
go work use ./drivers/mariadb
go work use ./drivers/mongodb
go work use ./drivers/mysql
go work use ./drivers/oracle
go work use ./drivers/postgresql
go work use ./drivers/redis
go work use ./drivers/sqlite
go work use ./drivers/sqlserver
go mod tidy
