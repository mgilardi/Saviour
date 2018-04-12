package core

// Cron takes in functions and adds them to a map that is run at a specified
// interval. It employs two go routines a controller/worker

// @TODO Pick an order for var, type, const and do the same in all files.

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
	options := OptionsHandler.GetOption("Core")
	cron.interval = time.Duration(int(options["Interval"].(float64))) * time.Minute
	cron.startInterval()
	CronHandler = &cron
}

// startInterval
// creates two goroutines one for running the time interval and initializing
// the Push function and one to allow the changing of the time interval. By providing the
// loaded interval when the cron has been completed.
func (cron *Cron) startInterval() {
	chanIntervalReset := make(chan bool)

	// "go func" creates a goroutine that gets the time interval from the interval controller thread
	// it will sleep the time interval provided by the controller then run the push of loaded cron
	// jobs when completed then send the interval reset signal to the interval controller thread
	go func() {
		for {
			time.Sleep(<-cron.intervalChan)
			Logger("CronStarting", "Cron", MSG)
			cron.Push()
			chanIntervalReset <- true
		}
	}()

	// "go func" Creates a new thread for the time interval controller loads the initial interval
	// the was loaded from options and sends the interval to the main cron thread. A for loop with a
	// select statement runs through the channels if it recieves a interval reset signal it will
	// return the currently loaded interval to the main cron thread, if it recieves a signal from the
	// new interal channel if will hold that as the current interval to be sent to the main cron
	// thread on the next completed cron thread
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

	// Starts the controller thread by sending the initial interval
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

// BatchWork will run all jobs located in an map of functions
func BatchWork(jobs map[int]func()) {
	Logger("Batch_Jobs_Recieved::"+strconv.Itoa(len(jobs)), "Cron", MSG)
	for _, job := range jobs {
		Work(job)
	}
}

// Work with run a worker on a go routine then waits for the job to be completed
func Work(job func()) {
	go func() {
		done := make(chan bool)
		go Run(done, job)
		<-done
		Logger("Job::Done", "Cron", MSG)
	}()
}

// Run the current job and return done on channel to the Work function
func Run(done chan<- bool, job func()) {
	job()
	done <- true
}
