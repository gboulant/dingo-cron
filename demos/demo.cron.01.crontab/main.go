package main

import (
	"fmt"

	cron "github.com/gboulant/dingo-cron"
)

func writedata() {
	fmt.Println("hello")
}

func main() {

	crontab := cron.NewCrontab()

	period := cron.D("5s")
	table := cron.NewPeriodicTable().Every(period).RoundedBy(cron.D("1s"))
	task := func() { writedata() }
	crontab.Add(task).WithTimetable(table)

	crontab.Start()
	crontab.Wait()
}
