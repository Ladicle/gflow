package gflow

import (
	"time"
)

// TaskRunner is interface based on all tasks.
type TaskRunner interface {
	Run() error
}

// NewPrintRunner generate PrintRUnner instance
func NewPrintRunner(message string) PrintRunner {
	return PrintRunner{message: message}
}

// PrintRunner implements TaskRunner
type PrintRunner struct {
	message string
}

// print message function
func (p PrintRunner) Run() error {
	time.Sleep(3 * time.Second)
	// TODO: logs print to file (!stdout)
	//fmt.Println(p.message)
	return nil
}

// HelmRunner is struct to do helm command
type HelmRunner struct{}

// comment
func (h HelmRunner) Run() error {
	return nil
}
