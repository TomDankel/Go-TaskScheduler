package main

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	EdfScheduler "github.com/edfScheduler"
	"github.com/task"
	"time"
)

var small int
var mid int
var high int

func main() {
	measruetime()
	edf := EdfScheduler.NewEdfScheduler()
	RunPeriodic(edf, 100, 200, "short period task", small)
	RunPeriodic(edf, 1000, 20, "middle period task", mid)
	RunPeriodic(edf, 10000, 2, "long period task", high)
	edf.Wg.Add(1)
	go edf.Run()
	time.Sleep(200 * time.Second)
	edf.EndScheduler()
	edf.Wg.Wait()
}

type Fibonacci struct {
	*task.TaskI
	iteration int
}

func NewFibonacci(iteration int) *Fibonacci {
	return &Fibonacci{task.NewTaskI(), iteration}
}

func (fib *Fibonacci) Run() {
	a := 0
	b := 1
	c := 0
	for i := 0; i < fib.iteration; i++ {
		exit := fib.Control()
		if exit {
			return
		}
		c = a
		a = b
		b = c + b
	}
	fib.Finished()
}

func RunPeriodic(edf *EdfScheduler.SchedulerI, period int, iteration int, name string, length int) {
	currentTime := time.Now().Add(time.Minute)
	for i := 0; i < iteration; i++ {
		fib := NewFibonacci(length)
		fib.SetName(name)
		var offset time.Duration
		offset = time.Duration(period * i)
		fib.SetDeadline(currentTime.Add(offset * time.Millisecond))
		edf.Schedule(fib)
	}
}

func measruetime() {
	a := 0
	b := 1
	c := 0
	deadline := time.Now().Add(10 * time.Millisecond)
	for i := 0; i < 100000000; i++ {
		if time.Now().After(deadline) {
			small = i + 2
			break
		}
		c = a
		a = b
		b = c + b
	}
	a = 0
	b = 1
	c = 0
	deadline = time.Now().Add(100 * time.Millisecond)
	for i := 0; i < 100000000; i++ {
		if time.Now().After(deadline) {
			mid = i + 2
			break
		}
		c = a
		a = b
		b = c + b
	}
	a = 0
	b = 1
	c = 0
	deadline = time.Now().Add(2000 * time.Millisecond)
	for i := 0; i < 100000000000; i++ {
		if time.Now().After(deadline) {
			high = i + 2
			break
		}
		c = a
		a = b
		b = c + b

	}
	fmt.Printf("Small: %d, mid: %d, high: %d\n", small, mid, high)
}

/*
	3 tasks:
		- 10ms execution alle 100ms
		- 100ms execution alle 1000ms
		- 2000ms execution alle 10000ms
	Ausgabe: name
*/
