package main

import (
	"fmt"
	"github.com/task"
	"time"
)

const N = 100

func main() {
	fmt.Println("Hello World")
	start := make(chan int)
	suspend := make(chan int)
	kill := make(chan int)
	threads := &task.TaskI{SuspendCh: start, PlayCh: suspend, KillCh: kill}
	multi := task.Task(threads)
	multi.Play(fibonacci)
	multi.Suspend()
	multi.Continue()
	time.Sleep(10)
	multi.Kill()
}

func fibonacci(task *task.TaskI) {
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
