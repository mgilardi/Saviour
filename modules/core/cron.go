package core

// Cron takes in functions and adds them to a map that is run at a specified
// interval. It employs two go routines a controller/worker

import (
	"fmt"
	"strconv"
	"time"
)

// CronHandler Global Interface
var CronHandler *Cron

// Cron hold channels, options, & jobs
type Cron struct {
	jobs         map[int]func()
	interval     time.Duration
	intervalChan chan time.Duration
	cronCount    int
}

// InitCron sets up the cronhandler
func InitCron() {
	var cron Cron
	cron.cronCount = 0
	cron.jobs = make(map[int]func())
	cron.intervalChan = make(chan time.Duration)
	cron.interval = time.Duration(10) * time.Second
	cron.startInterval()
	CronHandler = &cron
}

func (cron *Cron) startInterval() {
	go func() {
		for {
			time.Sleep(<-cron.intervalChan)
			cron.Push()
		}
	}()
	go func() {
		for {
			cron.intervalChan <- cron.interval
		}
	}()
}

// Add cron job
func (cron *Cron) Add(job func()) {
	cron.jobs[cron.cronCount] = job
	cron.cronCount++
}

// Push jobs to worker
func (cron *Cron) Push() {
	BatchWork(cron.jobs)
}

// Interval changes the time in between Cron activation
func (cron *Cron) Interval(interval int) {
	cron.interval = time.Duration(interval) * time.Second
}

// BatchWork will run all jobs located in an array of functions
func BatchWork(jobs map[int]func()) {
	fmt.Println("Batch_Jobs_Recieved::" + strconv.Itoa(len(jobs)))
	for _, job := range jobs {
		Work(job)
	}
}

// Work with run a worker on a go routine
func Work(job func()) {
	go func() {
		done := make(chan bool)
		go Run(done, job)
		<-done
		fmt.Println("Job::Done")
	}()
}

// Run the current job and return done on channel
func Run(done chan<- bool, job func()) {
	job()
	done <- true
}
