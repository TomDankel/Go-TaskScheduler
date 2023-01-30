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
	threads := &thread.Thread{SuspendCh: start, PlayCh: suspend, KillCh: kill}
	multi := thread.Multithreading(threads)
	multi.Suspend()
	fmt.Println(multi)
}
