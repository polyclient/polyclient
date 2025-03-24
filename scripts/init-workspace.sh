#!/usr/bin/env bash

rm -f go.work go.work.sum

go work init
go work use .
go work use ./drivers/sqlite
go work use ./drivers/postgres
go work use ./drivers/mysql
go work use ./drivers/mongodb
go mod tidy
