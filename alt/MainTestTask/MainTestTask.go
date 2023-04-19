package main

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	EdfScheduler "github.com/edfScheduler"
	"github.com/task"
	"time"
)

/*
basic concept of the scheduler:
-receives task (name of a function) and a deadline
-creates a task/thread(struct)
-runs the task in a GoRoutine scheduled via the deadlines (ascending - EDF)
-sorts the priority of the tasks new, if new task is scheduled
*/

const N = 100

func main() {
	//import EdfScheduler like above
	//instantiate new scheduler
	edf := EdfScheduler.NewEdfScheduler()
	//schedule a task (method/function) with corresponding deadline via calling edf.Schedule with method and deadline as parameter
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*7)) //each task gets an id from the scheduler starting with 0 (ascending) - ID: 0
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*6)) //ID: 1
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*8)) //ID: 2
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*4)) //ID: 3
	edf.Wg.Add(1)
	//start EdfScheduler in a separate GoRoutine
	go edf.Run()
	time.Sleep(time.Second * 1)
	//schedule new tasks while Scheduler is running
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*5)) //ID: 4
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*4)) //ID: 5
	time.Sleep(time.Second * 10)
	//end scheduler when finished
	edf.EndScheduler()
	edf.Wg.Wait()
}

func fibonacci(task task.Task) { //programmer/user needs to define his task (function/method) with the parameter task from type task.Task
	a := 0
	b := 1
	c := 0
	for i := 0; i < N; i++ {
		//task needs iterative task and programmer/user needs to call task.Control each iteration
		exit := task.Control()
		if exit == 1 {
			return
		}
		c = a
		a = b
		b = c + b
	}
	//when finished the programmer/user has to tell it to the Scheduler with task.Finished
	task.Finished()
}

type Fibonacci struct {
	*task.TaskI
}

func NewFibonacci() *Fibonacci {
	return &Fibonacci{task.NewTaskI()}
}

func (fib *Fibonacci) fibonacci() {
	a := 0
	b := 1
	c := 0
	for i := 0; i < N; i++ {
		exit := fib.Control()
		if exit == 1 {
			return
		}
		fmt.Println(a)
		c = a
		a = b
		b = c + b
	}
}
