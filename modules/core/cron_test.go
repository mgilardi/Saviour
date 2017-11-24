package core

import (
	"testing"
	"time"
)

func TestCron_InitCron(t *testing.T) {
	InitDebug(true)
	InitCron()
	CronHandler.Add(func() {
		Sys("Test", "Test")
	})
	CronHandler.Push()
	CronHandler.Interval(1)
	time.Sleep(50)
}
