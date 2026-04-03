package main

import (
	"fmt"
	"os"

	"github.com/aymenhmaidiwastaken/gitwise/cmd"
)

// Set via ldflags at build time.
var version = "dev"

func main() {
	cmd.Version = version
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
