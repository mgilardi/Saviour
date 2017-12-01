package core

import (
	"testing"
	"time"
)

func TestCron_InitCron(t *testing.T) {
	CronHandler.Add(func() {
		Logger("Test", "Test", MSG)
	})
	CronHandler.Push()
	CronHandler.Interval(1)
	time.Sleep(50)
}
