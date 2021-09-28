package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	LogActivate(false)

	m := map[string]Clock{
		"03:45:30": NewClock().H(3).M(45).S(30),
		"--:45:30": NewClock().M(45).S(30),
		"--:45:--": NewClock().M(45),
		"03:--:--": NewClock().H(3),
	}

	for ref, c := range m {
		res := c.String()
		if res != ref {
			t.Errorf("clock is %s (should be %s)", res, ref)
		}
	}

	for ref := range m {
		res := C(ref).String()
		if res != ref {
			t.Errorf("clock is %s (should be %s)", res, ref)
		}
	}

	LogActivate(true)
}

func duration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func buildSchedule(nsteps int, startdate time.Time, p time.Duration, c Clock, r time.Duration) []time.Time {
	logger.Printf("init: %s\n", time2string(startdate))
	timetable := NewPeriodicTable().Every(p).AtClock(c).RoundedBy(r)
	schedule := Schedule(timetable, startdate, nsteps)
	for _, timestamp := range schedule {
		logger.Printf("next: %s\n", time2string(timestamp))
	}
	return schedule
}

func checkSchedule(res []time.Time, ref []time.Time) error {
	for i := 0; i < len(ref); i++ {
		if res[i] != ref[i] {
			return fmt.Errorf("Date %d is %s (should be %s)", i, res[i], ref[i])
		}
	}
	return nil
}

func TestSchedule(t *testing.T) {
	LogActivate(false)

	// Note that this test function does not really test something.
	// It is more like an example of usage of the timetable and scheduler.

	ts, _ := time.Parse("2006 Jan 02 15:04:05", "2020 Nov 01 15:23:30")
	logger.Printf("init: %s\n", time2string(ts))

	log := func(message string) {
		logger.Println("--------------------------------------------------------------")
		logger.Println(message)
	}

	nsteps := 4

	// --------------------------------------------------------------
	log("Toutes les 30 minutes, arrondie à la dizaine de minutes près")
	clock := NewClock()
	period := duration("30m") // toutes les 30 minutes
	round := duration("10m")  // arondie à la dizaine de minutes
	schedule := buildSchedule(nsteps, ts, period, clock, round)
	err := checkSchedule(schedule, []time.Time{
		T("01 Nov 2020 15:50:00"),
		T("01 Nov 2020 16:20:00"),
		T("01 Nov 2020 16:50:00"),
		T("01 Nov 2020 17:20:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Toutes les 30 minutes, arrondie à la demi-heure près")
	clock = NewClock()
	period = duration("30m") // toutes les 30 minutes
	round = duration("30m")  // arondie à la demi-heure
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("01 Nov 2020 16:00:00"),
		T("01 Nov 2020 16:30:00"),
		T("01 Nov 2020 17:00:00"),
		T("01 Nov 2020 17:30:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Toutes les heures, et aux heures rondes")
	clock = NewClock()
	period = duration("1h") // toutes les heures les jours
	round = duration("1h")  // aux heures rondes
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("01 Nov 2020 16:00:00"),
		T("01 Nov 2020 17:00:00"),
		T("01 Nov 2020 18:00:00"),
		T("01 Nov 2020 19:00:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Tous les jours, arrondi à la journée près")
	clock = NewClock()
	period = duration("24h") // tous les jours
	round = duration("24h")  // arondie à la journée (déclenchement à minuit)
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("03 Nov 2020 00:00:00"),
		T("04 Nov 2020 00:00:00"),
		T("05 Nov 2020 00:00:00"),
		T("06 Nov 2020 00:00:00"),
	})
	if err != nil {
		t.Error(err)
	}
	// WARN: for this king of use case, it is better to use a clock("00:00:00")
	// because with a round of 24h we miss a complete day at the beginning

	// --------------------------------------------------------------
	log("Tous les jours, arrondi à la journée prés (déclenchement à minuit)")
	period = duration("24h") // tous les jours
	clock = C("00:00:00")    // déclenchement à minuit
	round = 0
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("02 Nov 2020 00:00:00"),
		T("03 Nov 2020 00:00:00"),
		T("04 Nov 2020 00:00:00"),
		T("05 Nov 2020 00:00:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Une fois par jour à une heure ronde autour de l'heure de démarrage")
	clock = NewClock()
	period = duration("24h") // tous les jours
	round = duration("1h")   // arondie à l'heure près ( (à partir de cet instant)
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("02 Nov 2020 15:00:00"),
		T("03 Nov 2020 15:00:00"),
		T("04 Nov 2020 15:00:00"),
		T("05 Nov 2020 15:00:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Une fois par jour à 13h46")
	period = duration("24h")       // tous les jours
	clock = NewClock().H(13).M(46) // déclecnchement à 13h46m
	round = 0
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("02 Nov 2020 13:46:00"),
		T("03 Nov 2020 13:46:00"),
		T("04 Nov 2020 13:46:00"),
		T("05 Nov 2020 13:46:00"),
	})
	if err != nil {
		t.Error(err)
	}

	// --------------------------------------------------------------
	log("Une fois par heure, à la minute 46 (par rapport à l'heure ronde)")
	period = duration("1h")  // tous les jours
	round = 0                // pas d'arrondi
	clock = NewClock().M(46) // à la minute 46
	schedule = buildSchedule(nsteps, ts, period, clock, round)
	err = checkSchedule(schedule, []time.Time{
		T("01 Nov 2020 16:46:00"),
		T("01 Nov 2020 17:46:00"),
		T("01 Nov 2020 18:46:00"),
		T("01 Nov 2020 19:46:00"),
	})
	if err != nil {
		t.Error(err)
	}

	LogActivate(true)
}
