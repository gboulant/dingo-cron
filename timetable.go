package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// -------------------------------------------------------------------------

const clockUnsetValue = -1

// Clock represents a time clock, i.e. a time of the day
type Clock interface {
	// H returns a new clock with the hour set to the specified digit
	H(digit int) Clock
	// M returns a new clock with the minute set to the specified digit
	M(digit int) Clock
	// S returns a new clock with the second set to the specified digit
	S(digit int) Clock
	// String returns a string representation of this clock like H:M:S
	String() string
}

type clock struct {
	hour   int
	minute int
	second int
}

// NewClock creates a new Clock instance
func NewClock() Clock {
	return clock{clockUnsetValue, clockUnsetValue, clockUnsetValue}
}

func (c clock) H(digit int) Clock {
	if digit >= 0 && digit < 24 {
		c.hour = digit
	}
	return c
}
func (c clock) M(digit int) Clock {
	if digit >= 0 && digit < 59 {
		c.minute = digit
	}
	return c
}
func (c clock) S(digit int) Clock {
	if digit >= 0 && digit < 59 {
		c.second = digit
	}
	return c
}

func (c clock) String() string {
	h := "--"
	m := "--"
	s := "--"
	if c.hour != clockUnsetValue {
		h = fmt.Sprintf("%.2d", c.hour)
	}
	if c.minute != clockUnsetValue {
		m = fmt.Sprintf("%.2d", c.minute)
	}
	if c.second != clockUnsetValue {
		s = fmt.Sprintf("%.2d", c.second)
	}
	return fmt.Sprintf("%s:%s:%s", h, m, s)
}

// -------------------------------------------------------------------------

func string2clock(value string) Clock {
	clock := clock{clockUnsetValue, clockUnsetValue, clockUnsetValue}
	tokens := strings.Split(value, ":")
	if len(tokens) != 3 {
		return clock
	}
	h, err := strconv.Atoi(tokens[0])
	if err == nil {
		clock.hour = h
	}
	m, err := strconv.Atoi(tokens[1])
	if err == nil {
		clock.minute = m
	}
	s, err := strconv.Atoi(tokens[2])
	if err == nil {
		clock.second = s
	}
	return clock
}

func string2time(value string) time.Time {
	datetime, _ := time.Parse("02 Jan 2006 15:04:05", value)
	return datetime
}
func time2string(t time.Time) string {
	return fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
}
func string2duration(value string) time.Duration {
	d, _ := time.ParseDuration(value)
	return d
}

// C is a rapid writing function to create a Clock from its string representation H:M:S
var C = string2clock

// T is a rapid writing function to create a timestamp from its string representation
var T = string2time

// S is a rapid writing function to create a string representation of a datetime
var S = time2string

// D is a rapid writing function to create a duration from its string representation
var D = string2duration

// -------------------------------------------------------------------------

// Timetable represents a schedule timetable. A timetable is able to give
// the next date of the schedule, from the given date.
type Timetable interface {
	Next(from time.Time) time.Time
}

// PeriodicTable is a Timetable with a periodic schedule
type PeriodicTable interface {
	Timetable
	Every(period time.Duration) PeriodicTable
	AtClock(c Clock) PeriodicTable
	RoundedBy(round time.Duration) PeriodicTable
}

type periodicTable struct {
	period time.Duration
	clock  clock
	round  time.Duration
	err    error
}

func (t *periodicTable) Every(period time.Duration) PeriodicTable {
	t.period = period
	if t.round > t.period {
		t.err = fmt.Errorf("You request a period duration smaller than the round duration")
	}
	return t
}

func (t *periodicTable) AtClock(c Clock) PeriodicTable {
	t.clock = c.(clock)
	return t
}

func (t *periodicTable) RoundedBy(round time.Duration) PeriodicTable {
	t.round = round
	if t.round > t.period {
		t.err = fmt.Errorf("You request a round duration greater than the period")
	}
	return t
}

// Next implements the interface Timetable
func (t periodicTable) Next(from time.Time) time.Time {
	timestamp := from.Add(t.period).Round(t.round)
	timestamp = timeAtClock(timestamp, t.clock)
	return timestamp
}

// NewPeriodicTable returns a void periodic timetable
func NewPeriodicTable() PeriodicTable {
	timetable := periodicTable{
		period: 0,
		clock:  NewClock().(clock),
		round:  0,
	}
	return &timetable
}

// timeAtClock returns a time rounds at the given clock in the day
func timeAtClock(t time.Time, c clock) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()
	loc := t.Location()

	if c.second != clockUnsetValue {
		sec = c.second
	}
	if c.minute != clockUnsetValue {
		min = c.minute
		if c.second == clockUnsetValue {
			sec = 0
		}
	}
	if c.hour != clockUnsetValue {
		hour = c.hour
		if c.minute == clockUnsetValue {
			min = 0
			sec = 0
		}
	}

	t = time.Date(year, month, day, hour, min, sec, 0, loc)
	return t
}

// Schedule returns the list of dates generated from the startdate using the specified timetable
func Schedule(timetable Timetable, startdate time.Time, nsteps int) []time.Time {
	var schedule []time.Time = make([]time.Time, nsteps)
	t := startdate
	for i := 0; i < nsteps; i++ {
		t = timetable.Next(t)
		schedule[i] = t
	}
	return schedule
}
