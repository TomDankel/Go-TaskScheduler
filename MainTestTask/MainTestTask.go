package main

import (
	EdfScheduler "github.com/edfScheduler"
	"time"
)

//authors: Tom Dankel and Luca Schwarz

func testTask() {
	small, mid, high := measuretime(10, 100, 1000)
	edf := EdfScheduler.NewEdfScheduler()
	RunPeriodic(edf, 100, 200, "short period task", small, 5)
	RunPeriodic(edf, 1000, 20, "middle period task", mid, 5)
	RunPeriodic(edf, 10000, 2, "long period task", high, 5)
	go edf.Run()
	time.Sleep(1 * time.Second)
	fib := NewFibonacci(high)
	fib.SetName("test long")
	fib.SetDeadline(time.Now().Add(2 * time.Second))
	edf.Schedule(fib)
	time.Sleep(30 * time.Second)
	edf.EndScheduler()
	edf.Wg.Wait()
}
