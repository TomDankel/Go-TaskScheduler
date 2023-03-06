package task

//authors: Tom Dankel and Luca Schwarz
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
	resumeCh  chan bool
	killCh    chan bool
}

func NewTaskI() *TaskI {
	task := &TaskI{}
	task.suspendCh = make(chan bool)
	task.resumeCh = make(chan bool)
	task.killCh = make(chan bool)
	return task
}

func (th *TaskI) PlayMethod(method func()) {
	go method()
}

func (th *TaskI) PlayFunction(method func(task Task)) {
	go method(th)
}
func (th *TaskI) Resume() {
	fmt.Println("Resume Go Routine")
	th.resumeCh <- true
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
		fmt.Println("Routine Killed")
		return 1
	}
	if th.checkChanel(th.suspendCh) {
		for true {
			fmt.Println("Routine Paused")
			if th.checkChanel(th.resumeCh) {
				fmt.Println("Routine Resumed")
				return 0
			}
			if th.checkChanel(th.killCh) {
				fmt.Println("Routine Killed")
				return 1
			}
		}
	}
	return 0
}

func (th *TaskI) checkChanel(chanel chan bool) bool {
	select {
	case message := <-chanel:
		if message {
			fmt.Printf("read from chanel: %t\n", message)
			return true
		}
		return false
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
