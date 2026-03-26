package main

import (
	"os"

	"github.com/go-tapd/cli/internal/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
