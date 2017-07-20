package gflow

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// TODO: convert to configuration value
const (
	maxWorkerNum       = 2
	taskTimeoutMinutes = 30
)

// NewWorkflow generate workflow instance.
func NewWorkflow(tasks ...*Task) Workflow {
	w := Workflow{tasks: []*Task{}}
	for _, t := range tasks {
		w.tasks = append(w.tasks, t)
	}
	return w
}

// Workflow is struct contains tasks
type Workflow struct {
	tasks []*Task
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
			log.Printf("start: %s", task.Name)
			if err := task.Execute(); err != nil {
				return err
			}
			log.Printf("finish: %s", task.Name)
		case <-time.After(taskTimeoutMinutes * time.Minute):
			return fmt.Errorf(
				"TimeoutError: Task execution time exceeded over %d minute.",
				taskTimeoutMinutes)
		}
	}
}
