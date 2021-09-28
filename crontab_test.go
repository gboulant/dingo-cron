package cron

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCrontabPrint(t *testing.T) {
	LogActivate(false)

	cron := crontab{}

	// Add a cron task in the cron table is quite simple
	date, _ := time.Parse("02 Jan 2006 15:04:05", "01 Nov 2020 15:23:30")
	table := NewPeriodicTable().Every(D("24h")).AtClock(C("13:47:00"))
	task := func() { fmt.Println("Hello, task 01 is working") }
	cron.Add(task).WithName("task01").WithTimetable(table).StartAfter(date)

	// You may also used a condensed writing
	cron.Add(func() { fmt.Println("Hello, task 02 is working") }).WithName("task02").
		WithTimetable(NewPeriodicTable().Every(D("1h")).AtClock(C("--:47:--"))).
		StartAfter(T("01 Nov 2020 15:23:30"))

	logger.Printf("crontab:\n%s\n", cron)

	LogActivate(false)
}

func TestCrontabExec(t *testing.T) {
	LogActivate(false)

	runtest := os.Getenv("RUN_LONG_TEST")
	if runtest != "1" {
		t.Skip("RUN_LONG_TEST != 1 ==> SKIP")
	}

	cron := crontab{}

	cron.Add(func() { logger.Println("Hello, task 01 is working") }).
		WithName("task01").
		WithTimetable(NewPeriodicTable().Every(D("3s"))).StartAfter(time.Now().Add(D("4s"))).
		NbTimes(2)

	cron.Add(func() { logger.Println("Hello, task 02 is working") }).
		WithName("task02").WithTimetable(NewPeriodicTable().Every(D("1s"))).
		NbTimes(9)

	cron.Start()
	cron.Wait()

	logger.Println("terminated")

	LogActivate(false)
}
