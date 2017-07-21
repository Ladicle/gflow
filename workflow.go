package gflow

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// TODO: convert to configuration value
const (
	maxWorkerNum       = 2
	taskTimeoutMinutes = 30
)

// NewWorkflow generate workflow instance.
func NewWorkflow(name string, tasks ...*Task) Workflow {
	w := Workflow{
		Name:  name,
		tasks: []*Task{},
	}
	for _, t := range tasks {
		w.tasks = append(w.tasks, t)
	}
	return w
}

// Workflow is struct contains tasks
type Workflow struct {
	Name  string
	tasks []*Task
}

// PrintTasks is function to print task dependencies
func (w Workflow) PrintTasks() {
	msg := []string{}
	for _, t := range w.tasks {
		depNames := []string{}
		for _, dt := range t.deps {
			depNames = append(depNames, dt.Name)
		}
		if len(depNames) == 0 {
			msg = append(msg, t.Name)
		} else {
			msg = append(msg, fmt.Sprintf("%s (bloked: %s)", t.Name, strings.Join(depNames, " & ")))
		}
	}
	fmt.Printf("%s\n\n", strings.Join(msg, "\n"))
}

// Add is function to add task
func (w *Workflow) Add(t *Task) {
	w.tasks = append(w.tasks, t)
}

// Start is function to start workflow
func (w Workflow) Start() error {
	w.tasks = NewDag(w.tasks).Tsort()
	c := make(chan Task, 2)
	eg := errgroup.Group{}

	for i := 0; i < maxWorkerNum; i++ {
		eg.Go(func() error { return runWorker(c) })
	}

	for _, t := range w.tasks {
		c <- *t
	}
	close(c)

	return eg.Wait()
}

func runWorker(c chan Task) error {
	for {
		select {
		case task, open := <-c:
			if !open {
				return nil
			}
			if err := task.Execute(); err != nil {
				return err
			}
		case <-time.After(taskTimeoutMinutes * time.Minute):
			return fmt.Errorf(
				"TimeoutError: Task execution time exceeded over %d minute.",
				taskTimeoutMinutes)
		}
	}
}
