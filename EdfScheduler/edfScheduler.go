package EdfScheduler

//authors: Tom Dankel and Luca Schwarz
import (
	"fmt"
	"github.com/task"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type EdfScheduler interface {
	Schedule(task.Task)
	insertToJobs(job)
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
	return scheduler
}

func (s *SchedulerI) Schedule(obj task.Task) {
	job := job{
		function: obj.Run,
		Deadline: obj.GetDeadline(),
		task:     obj,
		id:       s.id,
		run:      false,
	}
	s.insertToJobs(job)
	s.id++
}

func (s *SchedulerI) insertToJobs(job job) {
	inserted := false
	if len(s.jobs) == 0 {
		s.jobs = append(s.jobs, job)
	} else {
		for index, currentJob := range s.jobs {
			if job.Deadline.Before(currentJob.Deadline) {
				mutex.Lock()
				s.jobs = append(s.jobs[:index+1], s.jobs[index:]...)
				s.jobs[index] = job
				mutex.Unlock()
				inserted = true
				break
			}
		}
		if inserted == false {
			mutex.Lock()
			s.jobs = append(s.jobs, job)
			mutex.Unlock()
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
	iteration := false
	for {
		select {
		case abort := <-s.quit:
			if abort {
				fmt.Println("End EDF Scheduler")
				if (len(s.jobs) > 0) && (!currentJob.task.CheckFinished()) {
					//currentJob.task.Kill()
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
					if !removed {
						fmt.Println("paused 12433")
						fmt.Println(currentJob.id)
						currentJob.task.Suspend()
						fmt.Println("333")
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
				mutex.Lock()
				s.jobs = remove(s.jobs)
				mutex.Unlock()
				removed = true
				fmt.Println("Removed job")
				continue
			}
			if currentJob.task.CheckFinished() {
				fmt.Printf("finished %d with name: %s\n", currentJob.id, currentJob.task.GetName())
				mutex.Lock()
				s.jobs = remove(s.jobs)
				mutex.Unlock()
				removed = true
				continue
			}
			if currentJob.run {
				currentJob.task.Resume()
			} else {
				currentJob.task.PlayMethod(currentJob.function)
			}
			time.Sleep(time.Second)
		}
	}
}

func remove(slice []job) []job {
	slice = slice[1:]
	return slice
}
