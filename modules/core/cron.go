package core

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
	interval      time.Duration
	chanInterval  chan time.Duration
	chanLock      chan bool
	chanCheckLock chan bool
	chanRun       chan bool
	jobs          map[string]*CronObj
}

// InitCron this initializes the cron module
func InitCron() {
	var cron Cron
	chanJobs := make(chan map[string]*CronObj)
	cron.jobs = make(map[string]*CronObj)
	cron.chanInterval = make(chan time.Duration)
	cron.chanLock = make(chan bool)
	cron.chanCheckLock = make(chan bool)
	cron.chanRun = make(chan bool)
	options := GetOptions("Core")
	cron.interval = time.Duration(options["Interval"].(float64)) * time.Second
	go cron.startCron(chanJobs)
	go cron.startInterval(chanJobs)
	cron.chanInterval <- cron.interval
	CronHandler = &cron
	DebugHandler.Sys("Starting..", "Cron")
}

// startInterval is the controller for the cron worker. It controls when the
// worker starts by writing the jobs map to the channel, releasing it from
// blocking. Then it waits for the jobs map to be retrieved and overwrites
// cron jobs
func (cron *Cron) startInterval(chanJobs chan map[string]*CronObj) {
	interval := <-cron.chanInterval
	go func() {
		for {
			time.Sleep(<-cron.chanInterval)
			DebugHandler.Sys("Running..", "Cron")
			chanJobs <- cron.jobs
		}
	}()
	cron.chanInterval <- interval
	for {
		var locked bool
		locked = false
		select {
		case run := <-cron.chanRun:
			if run && !locked {
				DebugHandler.Sys("ForcedRun", "Cron")
				chanJobs <- cron.jobs
			}
		case interval = <-cron.chanInterval:
			DebugHandler.Sys("CronIntervalChanged", "Cron")
		case cron.jobs = <-chanJobs:
			DebugHandler.Sys("Finished..", "Cron")
			cron.chanInterval <- interval
		case lock := <-cron.chanLock:
			if lock {
				DebugHandler.Sys("CronLocked", "Cron")
				locked = true
			} else {
				DebugHandler.Sys("CronUnlocked", "Cron")
				locked = false
			}
		case isLocked := <-cron.chanCheckLock:
			if isLocked {
				cron.chanCheckLock <- locked
			}
		default:
			// Ignore
		}
	}
}

// startCron is the worker it will load the map of jobs to be done
// iterate through them and remove the ones that have the reapeat flag
// set to false. Then it sends the jobs map back to the controller
func (cron *Cron) startCron(chanJobs chan map[string]*CronObj) {
	for {
		jobs := <-chanJobs
		cron.chanLock <- true
		for key, job := range jobs {
			DebugHandler.Sys("LoadingJob::"+job.name, "Cron")
			job.run()
			if !job.repeat {
				delete(jobs, key)
			}
			LogHandler.Stat("Finished::"+job.name, "Cron")
		}
		chanJobs <- jobs
		cron.chanLock <- false
	}
}

// Interval changes cron interval
func (cron *Cron) Interval(interval int) {
	cron.chanInterval <- time.Duration(interval) * time.Second
}

// ForceStart forces cron to start
func (cron *Cron) ForceStart() {
	cron.chanRun <- true
}

// Add adds a new job to the jobs map if cronlock is enabled it will create a
// gothread the holds the new job until the lock is disengaged
func (cron *Cron) Add(cronName string, repeat bool, job func()) {
	newJob := newCronObj(cronName, repeat, job)
	cron.chanCheckLock <- true
	locked := <-cron.chanCheckLock
	if !locked {
		cron.jobs[cronName] = newJob
	} else {
		go func() {
			for {
				time.Sleep(time.Duration(10) * time.Second)
				cron.chanCheckLock <- true
				locked = <-cron.chanCheckLock
				if !locked {
					cron.jobs[cronName] = newJob
					break
				}
			}
		}()
	}
}
