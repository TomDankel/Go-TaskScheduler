package task

import "fmt"

type Task interface {
	PlayFunction(func(Task))
	PlayMethod(func())
	Resume()
	Suspend()
	Kill()
	Control() int
}

type TaskI struct {
	suspendCh chan bool
	playCh    chan bool
	killCh    chan bool
}

func NewTaskI() *TaskI {
	t := &TaskI{}
	t.suspendCh = make(chan bool)
	t.playCh = make(chan bool)
	t.killCh = make(chan bool)
	return t
}

func (th *TaskI) PlayMethod(method func()) {
	go method()
}

func (th *TaskI) PlayFunction(method func(task Task)) {
	go method(th)
}
func (th *TaskI) Resume() {
	fmt.Println("Resume Go Routine")
	th.playCh <- true
}
func (th *TaskI) Suspend() {
	fmt.Println("Suspend Go Routine")
	th.suspendCh <- true
}
func (th *TaskI) Kill() {
	fmt.Println("Kill Go Routine")
	th.killCh <- true
}
func (th *TaskI) Control() int {
	fmt.Println("enter controll")
	if th.checkChanel(th.killCh) {
		return 1
	}
	if th.checkChanel(th.suspendCh) {
		for true {
			fmt.Println("Routine Paused")
			if th.checkChanel(th.playCh) {
				return 0
			}
			if th.checkChanel(th.killCh) {
				return 1
			}
		}
	}
	return 0
}

func (th *TaskI) checkChanel(chanel chan bool) bool {
	select {
	case x, ok := <-chanel:
		if ok {
			fmt.Printf("read from chanel: %t", x)
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

/*
control structure in provided method (Recursion not supported?)
control function to pay attention of channels? --> not possible to communicate to method to pause
control function use inside provided method? approach a bit more complex than in other languages but possible
every x time units check for signal in channel --> in controll Pause or kill process (kill logic in provided method?)
	Pause and play possible in control method
*/
