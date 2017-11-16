package main

import "time"

var CronHandler *Cron

type Job func()

func (f Job) Run() { f() }

type CronObj struct {
	name, desc string
	cronChan   chan string
	reqJob     Job
	repeat     bool
}

func newCronObj(name string, desc string, repeat bool, job Job) *CronObj {
	var obj CronObj
	obj.name = name
	obj.desc = desc
	obj.repeat = repeat
	return &obj
}

type Cron struct {
	interval time.Duration
	options  map[string]interface{}
	jobs     map[string]*CronObj
}

func InitCron() {
	var newCron Cron
	newCron.options = GetOptions("Cron")
	newCron.interval = time.Duration(newCron.options["Interval"].(float64) * time.Minute.Minutes())
	newCron.jobs = make(map[string]*CronObj)
	CronHandler = &newCron
	go CronHandler.start()
}

func (cron *Cron) start() {
	for {
		count := 0
		key := make(chan string)
		value := make(chan *CronObj)
		done := make(chan bool)
		DebugHandler.Sys("CronStarting..", "Cron")
		time.Sleep(cron.interval)
		DebugHandler.Sys("CronRunning..", "Cron")
		for k, v := range cron.jobs {
			count++
			key <- k
			value <- v
			go func() {
				DebugHandler.Sys("RunningJob::"+<-key, "Cron")
				cron := <-value
				cron.reqJob.Run()
				done <- true
			}()
		}
	}
}

func (cron *Cron) Add(cronName string, desc string, repeat bool, job Job) {
	newJob := newCronObj(cronName, desc, repeat, job)
	cron.jobs[cronName] = newJob
}
