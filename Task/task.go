package task

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	"time"
)

type Task interface {
	PlayFunction(func(Task))
	PlayMethod(func())
	Resume()
	Suspend()
	Kill()
	Control() int
	Finished()
	CheckFinished() bool
	SetName(string)
	GetName() string
	SetDeadline(time.Time)
	GetDeadline() time.Time
	Run()
}

type TaskI struct {
	suspendCh  chan bool
	resumeCh   chan bool
	killCh     chan bool
	finishedCh chan bool
	name       string
	deadline   time.Time
}

func NewTaskI() *TaskI {
	task := &TaskI{}
	task.suspendCh = make(chan bool)
	task.resumeCh = make(chan bool)
	task.killCh = make(chan bool)
	task.finishedCh = make(chan bool)
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
func (th *TaskI) Finished() {
	fmt.Println("Finished Go Routine")
	th.finishedCh <- true
}
func (th *TaskI) Suspend() {
	fmt.Println("Suspend Go Routine")
	th.suspendCh <- true
}
func (th *TaskI) Kill() {
	fmt.Println("Kill Go Routine")
	th.killCh <- true
}
func (th *TaskI) SetName(name string) {
	th.name = name
}
func (th *TaskI) GetName() string {
	return th.name
}
func (th *TaskI) SetDeadline(deadline time.Time) {
	th.deadline = deadline
}
func (th *TaskI) GetDeadline() time.Time {
	return th.deadline
}
func (th *TaskI) Run() {

}
func (th *TaskI) Control() int {
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

func (th *TaskI) CheckFinished() bool {
	select {
	case message := <-th.finishedCh:
		if message {
			fmt.Println("finished checked")
			return true
		}
		return false
	default:
		return false
	}
}
