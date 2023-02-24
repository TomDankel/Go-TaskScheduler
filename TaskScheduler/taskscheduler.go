package main

import (
	"fmt"
	"github.com/thread"
)

func main() {
	fmt.Println("Hello World")
	start := make(chan int)
	suspend := make(chan int)
	kill := make(chan int)
	threads := &thread.TaskI{SuspendCh: start, PlayCh: suspend, KillCh: kill}
	multi := thread.Task(threads)
	multi.Suspend()
	fmt.Println(multi)
}
