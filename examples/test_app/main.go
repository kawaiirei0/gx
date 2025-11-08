package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Hello from gx cross-builder!\n")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
}
