package gflow

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

// NewTask is function to generate task instance.
func NewTask(name string, r TaskRunner, deps ...*Task) Task {
	t := Task{
		Name:     name,
		deps:     []*Task{},
		blockWgs: []*sync.WaitGroup{},
		runner:   r,
		spinner:  spinner.New(spinner.CharSets[9], 100*time.Millisecond),
	}
	if err := t.spinner.Color("yellow"); err != nil {
		log.Fatalln(err)
	}
	t.spinner.Suffix = ""

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
	spinner  *spinner.Spinner
}

// Ready is function to wait dependency tasks are completed.
func (t Task) ready() {
	if t.depWg != nil {
		t.SetMessage("wait dependency tasks")
		t.depWg.Wait()
	}
}

// Execute is function to run taskrunner and return result flag.
func (t Task) Execute() error {
	t.spinner.Start()
	t.ready()
	t.SetMessage("Execute task")
	if err := t.runner.Run(); err != nil {
		return err
	}
	t.spinner.Stop()
	for _, wg := range t.blockWgs {
		wg.Done()
	}
	return nil
}

// SetMessage is function to set spinner message
func (t Task) SetMessage(msg string) {
	t.spinner.Suffix = fmt.Sprintf("[%s] %s: ", t.Name, msg)
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
