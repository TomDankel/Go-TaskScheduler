package EdfScheduler

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	"github.com/task"
	"sync"
	"time"
)

var mutexJob = &sync.Mutex{}
var mutexId = &sync.Mutex{}

type EdfScheduler interface {
	Schedule(task.Task)
	Run()
	EndScheduler()
}

type SchedulerI struct {
	jobs []job
	quit chan bool
	id   int
	Wg   sync.WaitGroup
}

type job struct {
	function func()
	Deadline time.Time
	task     task.Task
	id       int
	run      bool
}

func NewEdfScheduler() *SchedulerI {
	scheduler := &SchedulerI{
		quit: make(chan bool),
		jobs: make([]job, 0),
		id:   0,
	}
	scheduler.Wg.Add(1)
	return scheduler
}

func (s *SchedulerI) Schedule(obj task.Task) {
	mutexId.Lock()
	job := job{
		function: obj.Run,
		Deadline: obj.GetDeadline(),
		task:     obj,
		id:       s.id,
		run:      false,
	}
	s.id++
	mutexId.Unlock()
	s.insertToJobs(job)
}

func (s *SchedulerI) insertToJobs(job job) {
	inserted := false
	if len(s.jobs) == 0 {
		s.jobs = append(s.jobs, job)
	} else {
		for index, currentJob := range s.jobs {
			if job.Deadline.Before(currentJob.Deadline) {
				mutexJob.Lock()
				s.jobs = append(s.jobs[:index+1], s.jobs[index:]...)
				s.jobs[index] = job
				mutexJob.Unlock()
				inserted = true
				break
			}
		}
		if inserted == false {
			mutexJob.Lock()
			s.jobs = append(s.jobs, job)
			mutexJob.Unlock()
		}
	}
}

func (s *SchedulerI) EndScheduler() {
	fmt.Println("End Scheduler")
	s.quit <- true
}

func (s *SchedulerI) Run() {
	defer s.Wg.Done()
	var currentJob job
	var removed bool
	switched := true
	iteration := false
	for {
		select {
		case abort := <-s.quit:
			if abort {
				fmt.Println("End EDF Scheduler")
				if (len(s.jobs) > 0) && (!currentJob.task.CheckFinished()) {
					currentJob.task.Kill()
				}
				return
			}
		default:
		}
		if len(s.jobs) > 0 {

			if !iteration {
				currentJob = s.jobs[0]
				iteration = true
			} else {
				if currentJob.id != s.jobs[0].id {
					switched = true
					if !removed {
						fmt.Println(currentJob.id)
						currentJob.task.Suspend()
						currentJob.run = true
						s.insertToJobs(currentJob)
					}
					currentJob = s.jobs[0]
					removed = false
				}
			}
			if currentJob.Deadline.Before(time.Now()) {
				fmt.Printf("Missed Deadline for job %d with name: %s\n", currentJob.id, currentJob.task.GetName())
				currentJob.task.Kill()
				mutexJob.Lock()
				s.jobs = remove(s.jobs)
				mutexJob.Unlock()
				removed = true
				fmt.Println("Removed job")
				continue
			}
			if currentJob.task.CheckFinished() {
				fmt.Printf("finished %d with name: %s\n", currentJob.id, currentJob.task.GetName())
				mutexJob.Lock()
				s.jobs = remove(s.jobs)
				mutexJob.Unlock()
				removed = true
				continue
			}
			if switched {
				switched = false
				if currentJob.run {
					currentJob.task.Resume()
				} else {
					currentJob.task.PlayMethod(currentJob.function)
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func remove(slice []job) []job {
	slice = slice[1:]
	return slice
}
