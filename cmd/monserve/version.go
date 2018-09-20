package main

import (
	"os"

	"github.com/glynternet/mon/internal/versioncmd"
)

// to be changed using ldflags with the go build command
var version = "unknown"

func init() {
	cmdDBServe.AddCommand(versioncmd.New(version, os.Stdout))
}
