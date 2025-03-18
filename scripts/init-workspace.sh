#!/usr/bin/env bash

rm -f go.work go.work.sum

go work init
go work use .
go work use ./plugins/sqlite
go mod tidy
