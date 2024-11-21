package main

import (
	. "github.com/cpunion/go-python"
)

func main() {
	Initialize()
	defer Finalize()
	println("Hello, World!")
}
