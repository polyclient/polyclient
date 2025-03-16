#!/usr/bin/env bash
set -e

# Clean existing bindings
rm -rf bindings/polyclient

# Install wit-bindgen-go
go get go.bytecodealliance.org/cmd/wit-bindgen-go

# Generate with verbose output
if go run go.bytecodealliance.org/cmd/wit-bindgen-go generate -v -o bindings/ ./wit; then
    echo "✅ Successfully generated Go bindings from WIT files"
else
    echo "❌ Failed to generate bindings"
    exit 1
fi
