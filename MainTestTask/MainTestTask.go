package main

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	EdfScheduler "github.com/edfScheduler"
	"github.com/task"
	"time"
)

const N = 100

func main() {
	edf := EdfScheduler.NewEdfScheduler()
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*3), time.Minute*3)
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*2), time.Minute*2)
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*4), time.Minute*4)
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*1), time.Minute*1)
	//go edf.Run()
	//time.Sleep(time.Second * 15)
	//fmt.Println("Hello PLayfunction")
	/*taskuse := task.NewTaskI()
	taskuse.PlayFunction(fibonacci)
	taskuse.Suspend()
	taskuse.Resume()
	time.Sleep(10)
	taskuse.Kill()

	// new edfsh
	//sh (fib 11:32
	//sh fib 11:30

	//end edf und dann alle kill
	//waitgroup für alle go routinen

	fmt.Println("Hello PlayMethod")
	fib := NewFibonacci()
	fib.PlayMethod(fib.fibonacci)
	fib.Suspend()
	fib.Resume()
	time.Sleep(10)
	fib.Kill()*/

}

func fibonacci(task task.Task) {
	a := 0
	b := 1
	c := 0
	for i := 0; i < N; i++ {
		exit := task.Control()
		if exit == 1 {
			return
		}
		fmt.Println(a)
		c = a
		a = b
		b = c + b
	}
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
