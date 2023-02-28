package task

import "fmt"

type Task interface {
	Continue()
	Suspend()
	Play(func(*TaskI))
	Kill()
	Control() int
}

type TaskI struct {
	SuspendCh chan bool
	PlayCh    chan bool
	KillCh    chan bool
}

func (th *TaskI) Play(method func(task *TaskI)) {
	go method(th)
}
func (th *TaskI) Continue() {
	fmt.Println("Continue Go Routine")
	th.PlayCh <- true
}
func (th *TaskI) Suspend() {
	fmt.Println("Suspend Go Routine")
	th.SuspendCh <- true
}
func (th *TaskI) Kill() {
	fmt.Println("Kill Go Routine")
	th.KillCh <- true
}
func (th *TaskI) Control() int {
	fmt.Println("enter controll")
	if <-th.KillCh {
		return 1
	}
	if <-th.SuspendCh {
		for true {
			fmt.Println("Routine Paused")
			if <-th.PlayCh {
				return 0
			}
			if <-th.KillCh {
				return 1
			}
		}
	}
	return 0
}

/*
control structure in provided method (Recursion not supported?)
control function to pay attention of channels? --> not possible to communicate to method to pause
control function use inside provided method? approach a bit more complex than in other languages but possible
every x time units check for signal in channel --> in controll Pause or kill process (kill logic in provided method?)
	Pause and play possible in control method
*/
