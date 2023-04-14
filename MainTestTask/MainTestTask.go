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
	fib1 := NewFibonacci(small)
	fib1.SetName("short period task")
	fib2 := NewFibonacci(mid)
	fib2.SetName("middle period task")
	fib3 := NewFibonacci(high)
	fib3.SetName("long period task")
	fib1.SetDeadline(time.Now().Add(time.Minute * 3))
	fib2.SetDeadline(time.Now().Add(time.Minute * 2))
	fib3.SetDeadline(time.Now().Add(time.Minute * 1))
	fib1.RunPeriodic(edf, 100, 200)
	fib2.RunPeriodic(edf, 1000, 20)
	fib3.RunPeriodic(edf, 10000, 2)
	edf.Wg.Add(1)
	go edf.Run()
	time.Sleep(10 * time.Second)
	edf.EndScheduler()
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
		if exit == 1 {
			return
		}
		c = a
		a = b
		b = c + b
	}
	fib.Finished()
}

func (fib *Fibonacci) RunPeriodic(edf *EdfScheduler.SchedulerI, period int, iteration int) {
	currentTime := time.Now().Add(time.Minute)
	for i := 0; i < iteration; i++ {
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
