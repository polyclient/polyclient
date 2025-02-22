package main

import (
	extism "github.com/extism/go-pdk"
)

//go:wasmexport greet
func greet() int32 {
	input := extism.Input()
	greeting := "Hello, " + string(input) + "!"
	extism.OutputString(greeting)
	return 0
}

func main() {}
