#!/usr/bin/env bash

if [ ! -f go.work ]; then
    go work init
    go work use .
    go mod tidy
fi
