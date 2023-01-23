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
	thread2 := thread.Thread{SuspendCh: start, PlayCh: suspend, KillCh: kill}
	thread2.Suspend()
}

//TODO extract Thread struct and methods to separate Class
/*
func (th *Thread) Play(method func(suspend chan int, play chan int, kill chan int)) {
	go method(th.suspend, th.play, th.kill)
}
func (th *Thread) Restart() {
	th.play <- 1
}
func (th *Thread) Suspend() {
	th.suspend <- 1
}
func (th *Thread) Kill() {
	th.kill <- 1
}
func (th *Thread) control() {

}
*/
/*
control structure in provided method (Recursion not supported?)
control function to pay attention of channels? --> not possible to communicate to method to pause
control function use inside provided method? approach a bit more complex than in other languages but possible
every x time units check for signal in channel --> in controll Pause or kill process (kill logic in provided method?)
	Pause and play possible in control method
*/
