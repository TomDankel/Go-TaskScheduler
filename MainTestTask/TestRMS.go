package main

import (
	EdfScheduler "github.com/edfScheduler"
	"time"
)

func testRMS() {
	small, mid, _ := measuretime(1, 3, 10)
	edf := EdfScheduler.NewEdfScheduler()
	RunPeriodic(edf, 6, 100, "A", mid, 1)
	RunPeriodic(edf, 8, 100, "B", mid, 1)
	RunPeriodic(edf, 10, 100, "C", small, 1)
	go edf.Run()
	time.Sleep(30 * time.Second)
	edf.EndScheduler()
	edf.Wg.Wait()
}
