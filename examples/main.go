package main

import (
	"log"

	"github.com/Ladicle/gflow"
)

func main() {
	t5 := gflow.NewTask("task5", gflow.NewPrintRunner("task5"))
	t4 := gflow.NewTask("task4", gflow.NewPrintRunner("task4"), &t5)
	t2 := gflow.NewTask("task2", gflow.NewPrintRunner("task2"))
	t3 := gflow.NewTask("task3", gflow.NewPrintRunner("task3"), &t5, &t4)
	t1 := gflow.NewTask("task1", gflow.NewPrintRunner("task1"), &t2, &t3)

	w := gflow.NewWorkflow("Print workflow", &t1, &t2, &t3, &t4, &t5)
	w.PrintTasks()
	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
