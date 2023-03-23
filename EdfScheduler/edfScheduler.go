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
	task *task.TaskI
	jobs []job
	quit chan bool
}

type job struct {
	function func(task task.Task)
	Deadline time.Time
	Duration time.Duration
}

func NewEdfScheduler() *SchedulerI {
	scheduler := &SchedulerI{}
	scheduler.task = task.NewTaskI()
	return scheduler
}

func (s *SchedulerI) Schedule(method func(task task.Task), deadline time.Time, duration time.Duration) {
	job := job{
		function: method,
		Deadline: deadline,
		Duration: duration,
	}
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

}
