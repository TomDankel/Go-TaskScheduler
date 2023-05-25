package main

import (
	EdfScheduler "github.com/edfScheduler"
	"time"
)

func testOverload() {
	small, mid, high := measuretime()
	edf := EdfScheduler.NewEdfScheduler()
	RunPeriodic(edf, 10, 200, "short period task", small, 0)
	RunPeriodic(edf, 100, 20, "middle period task", mid, 0)
	RunPeriodic(edf, 1000, 2, "long period task", high, 0)
	go edf.Run()
	time.Sleep(30 * time.Second)
	edf.EndScheduler()
	edf.Wg.Wait()
}
