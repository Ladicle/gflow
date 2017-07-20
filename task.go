package gflow

import (
	"log"
	"sync"
)

// TaskRunner is interface based on all tasks.
type TaskRunner interface {
	Run() error
}

// NewTask is function to generate task instance.
func NewTask(name string, r TaskRunner, deps ...*Task) Task {
	t := Task{
		Name:     name,
		deps:     []*Task{},
		blockWgs: []*sync.WaitGroup{},
		runner:   r,
	}
	for _, d := range deps {
		t.AddDep(d)
	}
	return t
}

// Task does not have dependency tasks.
// You need to implement Run method.
type Task struct {
	Name     string
	runner   TaskRunner
	blockWgs []*sync.WaitGroup
	depWg    *sync.WaitGroup
	deps     []*Task
}

// Ready is function to wait dependency tasks are completed.
func (t Task) ready() {
	if t.depWg != nil {
		log.Printf("%s wait dependency tasks", t.Name)
		t.depWg.Wait()
	}
}

// Execute is function to run taskrunner and return result flag.
func (t Task) Execute() error {
	t.ready()
	if err := t.runner.Run(); err != nil {
		return err
	}
	for _, wg := range t.blockWgs {
		wg.Done()
	}
	return nil
}

// AddDep is function to add dependency tasks.
func (t *Task) AddDep(dep *Task) {
	if t.depWg == nil {
		t.depWg = &sync.WaitGroup{}
	}
	t.depWg.Add(1)
	dep.blockWgs = append(dep.blockWgs, t.depWg)
	t.deps = append(t.deps, dep)
}
