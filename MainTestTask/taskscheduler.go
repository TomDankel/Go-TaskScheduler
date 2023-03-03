package main

import (
	"fmt"
	"github.com/task"
	"time"
)

const N = 100

func main() {
	fmt.Println("Hello PLayfunction")
	taskuse := task.NewTaskI()
	taskuse.PlayFunction(fibonacci)
	taskuse.Suspend()
	taskuse.Resume()
	time.Sleep(10)
	taskuse.Kill()

	fmt.Println("Hello PlayMethod")
	fib := NewFibonacci()
	fib.PlayMethod(fib.fibonacci)
	fib.Suspend()
	fib.Resume()
	time.Sleep(10)
	fib.Kill()

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
