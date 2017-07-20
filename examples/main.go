package main

import (
	"fmt"
	"log"

	"github.com/Ladicle/gflow"
)

// print message runner
type printRunner struct {
	message string
}

// print message function
func (p printRunner) Run() error {
	fmt.Println(p.message)
	return nil
}

func main() {
	t5 := gflow.NewTask("task5", printRunner{message: "task5"})
	t4 := gflow.NewTask("task4", printRunner{message: "task4"}, &t5)
	t2 := gflow.NewTask("task2", printRunner{message: "task2"})
	t3 := gflow.NewTask("task3", printRunner{message: "task3"}, &t5, &t4)
	t1 := gflow.NewTask("task1", printRunner{message: "task1"}, &t2, &t3)

	w := gflow.NewWorkflow(&t1, &t2, &t3, &t4, &t5)
	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
