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
	Schedule(func(task.Task), time.Time)
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
	function func(task task.Task)
	Deadline time.Time
	task     *task.TaskI
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

func (s *SchedulerI) Schedule(method func(task task.Task), deadline time.Time) {
	job := job{
		function: method,
		Deadline: deadline,
		task:     task.NewTaskI(),
		id:       s.id,
		run:      false,
	}
	s.id++
	s.insertToJobs(job)
}

func (s *SchedulerI) insertToJobs(job job) {
	mutex.Lock()
	inserted := false
	if len(s.jobs) == 0 {
		s.jobs = append(s.jobs, job)
	} else {
		for index, currentJob := range s.jobs {
			if job.Deadline.Before(currentJob.Deadline) {
				s.jobs = append(s.jobs[:index+1], s.jobs[index:]...)
				s.jobs[index] = job
				inserted = true
				break
			}
		}
		if inserted == false {
			s.jobs = append(s.jobs, job)
		}
	}
	mutex.Unlock()
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
				if len(s.jobs) > 0 {
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
					if !removed {
						currentJob.task.Suspend() //wenn eine neue kommt - deadlock
						currentJob.run = true
						s.insertToJobs(currentJob)
					}
					currentJob = s.jobs[0]
					removed = false
				}
			}
			if currentJob.Deadline.Before(time.Now()) {
				fmt.Printf("Missed deadline for job: %d", currentJob.id)
				currentJob.task.Kill()
				s.jobs = remove(s.jobs)
				removed = true
				continue
			}
			if currentJob.task.CheckFinished() {
				fmt.Print("with ID:")
				fmt.Println(s.jobs[0].id)
				s.jobs = remove(s.jobs)
				removed = true
				continue
			}
			if currentJob.run {
				currentJob.task.Resume()
			} else {
				currentJob.task.PlayFunction(currentJob.function)
			}
			time.Sleep(time.Second)
		}
	}
}

func remove(slice []job) []job {
	mutex.Lock()
	slice = slice[1:]
	mutex.Unlock()
	return slice
}
