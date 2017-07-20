package main

import (
	"fmt"
	"log"

	"github.com/Ladicle/gflow"
)

// print message
type printRunner struct {
	message string
}

// print message
func (p printRunner) Run() error {
	fmt.Println(p.message)
	return nil
}

func main() {
	t3 := gflow.NewTask("task3", printRunner{message: "task3"})
	t2 := gflow.NewTask("task2", printRunner{message: "task2"})
	t1 := gflow.NewTask("task1", printRunner{message: "task1"}, &t2)

	w := gflow.NewWorkflow()
	w.Add(t1)
	w.Add(t2)
	w.Add(t3)

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
