// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/extism/go-pdk"
	extism "github.com/extism/go-pdk"
)

//go:wasmexport query
func query() int32 {
	q := extism.InputString()
	pdk.Log(pdk.LogInfo, fmt.Sprintf("Executing query: %s", q))

	mockData := []struct {
		Name     string
		Email    string
		Birthday time.Time
	}{
		{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Birthday: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Jane Doe",
			Email:    "janedoe@example.com",
			Birthday: time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:     "Bob Smith",
			Email:    "bobsmith@example.com",
			Birthday: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	extism.OutputJSON(mockData)
	return 0
}

func main() {}
