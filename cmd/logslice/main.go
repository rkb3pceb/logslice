// Command logslice extracts time-range segments from large structured log files.
package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}
