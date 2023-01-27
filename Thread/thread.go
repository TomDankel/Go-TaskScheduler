package thread

type Multithreading interface {
	Continue()
	Suspend()
	Play(func(chan int, chan int, chan int))
	Kill()
	Control() int
}

type Thread struct {
	SuspendCh chan int
	PlayCh    chan int
	KillCh    chan int
}

func (th *Thread) Play(method func(suspend chan int, play chan int, kill chan int)) {
	go method(th.SuspendCh, th.PlayCh, th.KillCh)
}
func (th *Thread) Continue() {
	th.PlayCh <- 1
}
func (th *Thread) Suspend() {
	th.SuspendCh <- 1
}
func (th *Thread) Kill() {
	th.KillCh <- 1
}
func (th *Thread) Control() int {
	pause := <-th.SuspendCh
	kill := <-th.KillCh
	var play int
	if kill == 1 {
		return 1
	}
	if pause == 1 {
		for true {
			play = <-th.PlayCh
			if play == 1 {
				return 0
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
