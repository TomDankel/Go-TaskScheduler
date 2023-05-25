package main

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	EdfScheduler "github.com/edfScheduler"
	"github.com/task"
	"time"
)

func testTask() {
	small, mid, high := measuretime()
	edf := EdfScheduler.NewEdfScheduler()
	RunPeriodic(edf, 100, 200, "short period task", small, 5)
	RunPeriodic(edf, 1000, 20, "middle period task", mid, 5)
	RunPeriodic(edf, 10000, 2, "long period task", high, 5)
	go edf.Run()
	time.Sleep(20 * time.Second)
	fib := NewFibonacci(high)
	fib.SetName("test long")
	fib.SetDeadline(time.Now().Add(2 * time.Second))
	edf.Schedule(fib)
	time.Sleep(30 * time.Second)
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

func RunPeriodic(edf *EdfScheduler.SchedulerI, period time.Duration, iteration int, name string, length int, timeoffset time.Duration) {
	currentTime := time.Now().Add(timeoffset * time.Minute)
	for i := 0; i < iteration; i++ {
		fib := NewFibonacci(length)
		fib.SetName(name)
		var offset time.Duration
		offset = period * time.Duration(i)
		fib.SetDeadline(currentTime.Add(offset * time.Millisecond))
		edf.Schedule(fib)
	}
}

func measuretime() (int, int, int) {
	var small int
	var mid int
	var high int
	a := 0
	b := 1
	c := 0
	deadline := time.Now().Add(8 * time.Millisecond)
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
	deadline = time.Now().Add(90 * time.Millisecond)
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
	deadline = time.Now().Add(1900 * time.Millisecond)
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
	return small, mid, high
}
