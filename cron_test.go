package main

import (
	"testing"
	"time"
)

func TestCron_InitCron(t *testing.T) {
	InitDebug(true)
	InitCron()
	CronHandler.Add("Test", false, func() {
		DebugHandler.Sys("Test", "Test")
	})
	CronHandler.ForceStart()
	CronHandler.Interval(1)
	time.Sleep(50)
}
