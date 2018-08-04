package prl

import (
	"log"
	"os"
)

var (
	logger *log.Logger

	// Name of thread
	Name string

	// Version of application
	Version string
)

// Init the process
func Init(name, version string) *Task {
	logger = log.New(os.Stdout, "[parallel] ", 0)
	Name, Version = name, version
	logger.Printf("Initiating Thread [%s] with Version %s", name, version)
	return &Task{}
}
