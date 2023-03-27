package EdfScheduler

import (
	"fmt"
	"github.com/task"
	"time"
)

type EdfScheduler interface {
	Schedule(func(task.Task), time.Time, time.Duration)
	insertToJobs(job)
	Run()
	EndScheduler()
}

type SchedulerI struct {
	jobs []job
	quit chan bool
	id   int
}

type job struct {
	function func(task task.Task)
	Deadline time.Time
	Duration time.Duration
	task     *task.TaskI
	id       int
	run      bool
}

func NewEdfScheduler() *SchedulerI {
	scheduler := &SchedulerI{
		quit: make(chan bool),
		jobs: make([]job, 1),
		id:   0,
	}
	return scheduler
}

func (s *SchedulerI) Schedule(method func(task task.Task), deadline time.Time, duration time.Duration) {
	job := job{
		function: method,
		Deadline: deadline,
		Duration: duration,
		task:     task.NewTaskI(),
		id:       s.id,
		run:      false,
	}
	s.id++
	s.insertToJobs(job)
}

func (s *SchedulerI) insertToJobs(job job) {
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
}

func (s *SchedulerI) EndScheduler() {
	fmt.Println("End Scheduler")
	s.quit <- true
}

func (s *SchedulerI) Run() {
	var currentJob job
	iteration := false
	for {
		select {
		case abort := <-s.quit:
			if abort {
				fmt.Println("End EDF Scheduler")
				return
			}
		default:
		}
		if !iteration {
			currentJob = s.jobs[0]
			iteration = true
		} else {
			if currentJob.id != s.jobs[0].id {
				currentJob.task.Suspend()
				currentJob.run = true
				s.insertToJobs(currentJob)
				currentJob = s.jobs[0]
			}
		}
		if currentJob.Deadline.Before(time.Now()) {
			fmt.Printf("Missed deadline for job: %d", currentJob.id)
			s.jobs = remove(s.jobs)
		}
		if currentJob.task.CheckFinished() {
			s.jobs = remove(s.jobs)
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

func remove(slice []job) []job {
	slice = slice[1:]
	return slice
}
