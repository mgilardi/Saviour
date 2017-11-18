package main

// Cron takes in functions and adds them to a map that is run at a specified
// interval. It employs two go routines a controller/worker

import (
	"time"
)

// CronHandler Cron Global Variable
var CronHandler *Cron

// CronObj structure for cron job objects
type CronObj struct {
	name   string
	run    func()
	repeat bool
}

// newCronObj creates a new cron job object
func newCronObj(name string, repeat bool, job func()) *CronObj {
	var obj CronObj
	obj.name = name
	obj.run = job
	obj.repeat = repeat
	return &obj
}

// Cron is the main structure contains the interval and the jobs map
type Cron struct {
	interval time.Duration
	jobs     map[string]*CronObj
}

// InitCron this initializes the cron module
func InitCron() {
	var cron Cron
	chanJobs := make(chan map[string]*CronObj)
	cron.jobs = make(map[string]*CronObj)
	chanInterval := make(chan time.Duration)
	options := GetOptions("Cron")
	cron.interval = time.Duration(options["Interval"].(float64)) * time.Minute
	go cron.startCron(chanJobs)
	go cron.startInterval(chanJobs, chanInterval)
	chanInterval <- cron.interval
	CronHandler = &cron
	DebugHandler.Sys("Starting..", "Cron")
}

// startInterval is the controller for the cron worker. It controls when the
// worker starts by writing the jobs map to the channel, releasing it from
// blocking. Then it waits for the jobs map to be retrieved and overwrites
// cron jobs
func (cron *Cron) startInterval(chanJobs chan map[string]*CronObj, chanInterval chan time.Duration) {
	interval := <-chanInterval
	for {
		time.Sleep(interval)
		chanJobs <- cron.jobs
		cron.jobs = <-chanJobs
	}
}

// startCron is the worker it will load the map of jobs to be done
// iterate through them and remove the ones that have the reapeat flag
// set to false. Then it sends the jobs map back to the controller
func (cron *Cron) startCron(chanJobs chan map[string]*CronObj) {
	for {
		jobs := <-chanJobs
		for key, job := range jobs {
			DebugHandler.Sys("LoadingJob::"+job.name, "Cron")
			job.run()
			if !job.repeat {
				delete(jobs, key)
			}
		}
		chanJobs <- jobs
	}
}

// Add adds a new job to the jobs map
func (cron *Cron) Add(cronName string, repeat bool, job func()) {
	newJob := newCronObj(cronName, repeat, job)
	cron.jobs[cronName] = newJob
}
