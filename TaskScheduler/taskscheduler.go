package main

import (
	"fmt"
	"github.com/task"
)

func main() {
	fmt.Println("Hello World")
	start := make(chan int)
	suspend := make(chan int)
	kill := make(chan int)
	threads := &task.TaskI{SuspendCh: start, PlayCh: suspend, KillCh: kill}
	multi := task.Task(threads)
	multi.Suspend()
	fmt.Println(multi)
}
