package core

// Cron takes in functions and adds them to a map that is run at a specified
// interval. It employs two go routines a controller/worker

import (
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
	intervalSet  chan time.Duration
	cronCount    int
}

// InitCron sets up the cronhandler
func InitCron() {
	var cron Cron
	cron.cronCount = 0
	cron.jobs = make(map[int]func())
	cron.intervalChan = make(chan time.Duration)
	cron.intervalSet = make(chan time.Duration)
	options := OptionsHandler.GetOptions("Core")
	cron.interval = time.Duration(int(options["Interval"].(float64))) * time.Minute
	cron.startInterval()
	CronHandler = &cron
}

func (cron *Cron) startInterval() {
	chanIntervalReset := make(chan bool)
	go func() {
		for {
			time.Sleep(<-cron.intervalChan)
			Logger("CronStarting", "Cron", MSG)
			cron.Push()
			chanIntervalReset <- true
		}
	}()
	go func() {
		interval := <-cron.intervalSet
		cron.intervalChan <- interval
		for {
			select {
			case <-chanIntervalReset:
				cron.intervalChan <- interval
			case newInterval := <-cron.intervalSet:
				interval = newInterval
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
	cron.intervalSet <- cron.interval
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
	cron.intervalSet <- time.Duration(interval) * time.Second
}

// BatchWork will run all jobs located in an array of functions
func BatchWork(jobs map[int]func()) {
	Logger("Batch_Jobs_Recieved::"+strconv.Itoa(len(jobs)), "Cron", MSG)
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
		Logger("Job::Done", "Cron", MSG)
	}()
}

// Run the current job and return done on channel
func Run(done chan<- bool, job func()) {
	job()
	done <- true
}
