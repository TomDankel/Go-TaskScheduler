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
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*7))
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*6))
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*8))
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*4))
	edf.Wg.Add(1)
	go edf.Run()
	time.Sleep(time.Second * 1)
	edf.Schedule(fibonacci, time.Now().Add(time.Minute*5))
	time.Sleep(time.Second * 10)
	edf.EndScheduler()
	edf.Wg.Wait()
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
