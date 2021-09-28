package cron

import (
	"fmt"
	"hash/crc32"
	"time"
)

// Cronjob represents a cron job, i.e. a task executed repeatidly
// with a given timetable
type Cronjob interface {
	WithTimetable(table Timetable) Cronjob
	WithName(name string) Cronjob
	NbTimes(n int) Cronjob
	StartAfter(date time.Time) Cronjob
	String() string
	Start()
}

// Crontab represents a table of cron jobs, where a cron job is
// basically a task executed repeatidly with a given timetable
type Crontab interface {
	Add(task func()) Cronjob
	Start()
	Wait()
	String() string
}

// NewCrontab create a new Crontab
func NewCrontab() Crontab {
	return &crontab{}
}

// ------------------------------------------------------------

type cronjob struct {
	name       string
	task       func()
	timetable  Timetable
	startafter time.Time
	nbtimes    int

	ended chan time.Time
}

func (job *cronjob) WithTimetable(table Timetable) Cronjob {
	job.timetable = table
	return job
}

func (job *cronjob) WithName(name string) Cronjob {
	job.name = name
	return job
}

func (job *cronjob) NbTimes(n int) Cronjob {
	job.nbtimes = n
	return job
}

func (job *cronjob) StartAfter(date time.Time) Cronjob {
	job.startafter = date
	return job
}

func (job cronjob) String() string {
	s := fmt.Sprintf("name: %s\n", job.name)
	s += fmt.Sprintf("init: %s\n", time2string(job.startafter))
	schedule := Schedule(job.timetable, job.startafter, 5)
	for _, timestamp := range schedule {
		s += fmt.Sprintf("next: %s\n", time2string(timestamp))
	}
	return s[:len(s)-1]
}

func (job *cronjob) Start() {
	//logger.Printf("init[job:%s] begin\n", job.name)
	deadline := job.startafter
	job.submit(0, deadline)
	//logger.Printf("init[job:%s] end\n", job.name)
}

func (job *cronjob) submit(istep int, deadline time.Time) {
	//logger.Printf("submit[job:%s] step:%d begin\n", job.name, istep)

	// We first start a go routine that wait the timer fire
	if deadline.Before(time.Now()) {
		deadline = job.timetable.Next(time.Now())
	}
	taskended := make(chan time.Time, 1)
	timer := time.NewTimer(time.Until(deadline))
	go func() {
		<-timer.C
		logger.Printf("submit[job:%s,iter:%d] task: begin\n", job.name, istep)
		job.task()
		logger.Printf("submit[job:%s,iter:%d] task: end\n", job.name, istep)
		taskended <- time.Now()
	}()

	// Then we check if a next step is required
	if job.nbtimes > 0 && istep == job.nbtimes {
		logger.Printf("submit[job:%s] last submit\n", job.name)
		job.ended <- time.Now()
		return
	}

	// Finally, we start a second go routine that re-submits the job
	// when the current job is finished (it checks the channel taskended)
	go func() {
		<-taskended
		deadline = job.timetable.Next(deadline)
		logger.Printf("submit[job:%s,iter:%d] next run at %s\n", job.name, istep, deadline)
		job.submit(istep+1, deadline)
	}()
	//logger.Printf("submit[job:%s] step:%d end\n", job.name, istep)
}

func (job *cronjob) wait() {
	endedAt := <-job.ended
	logger.Printf("cron %s ended at %s\n", job.name, S(endedAt))
}

// ------------------------------------------------------------

type crontab struct {
	cronjobs []*cronjob
}

// hashInt returns an integer hash representation (hash.crc32) of the given string
func hashint(s string) uint32 {
	b := []byte(s)
	h := crc32.ChecksumIEEE(b)
	return h
}

func (c *crontab) Add(task func()) Cronjob {
	job := cronjob{task: task}
	job.name = fmt.Sprint(hashint(fmt.Sprint(&task)))
	job.ended = make(chan time.Time, 1)
	job.startafter = time.Now()
	c.cronjobs = append(c.cronjobs, &job)
	return &job
}

func (c crontab) Start() {
	for _, job := range c.cronjobs {
		job.Start()
	}
}

func (c crontab) Wait() {
	for _, job := range c.cronjobs {
		job.wait()
	}

}

func (c crontab) String() string {
	s := ""
	for i, job := range c.cronjobs {
		s += fmt.Sprintf("=== cron %.3d ===\n%s\n", i, job)
	}
	return s[:len(s)-1]
}
