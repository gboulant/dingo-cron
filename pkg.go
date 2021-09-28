// Package cron provides a simple function scheduler that works like a crontab.
package cron

import (
	"io/ioutil"
	"log"
	"os"
)

/*

Use cases
---------

UC01: Every day (period=24h), at 5h12m (clock=5h12m)
If Now is Mo-15h23m, the schedule is: Tu-5h12, We-5h12, Th-5h12, etc.

UC02: Every hour (period=60m), round to 10m (start at next hour+10m)
If Now is Mo-15h23m, the schedule is: Mo-16h10m, Mo-17h10m, Mo-17h10m, etc

UC03: Every half hour (period=30m), round to 10m (start at next hour+10m)
If Now is Mo-15h23m, the schedule is: Mo-16h10m, Mo-16h40m, Mo-17h10m, etc

UC04: Every hour (period=60m), no other indication
If Now is Mo-15h23m, the schedule is: Mo-15h23m, Mo-16h23m, Mo-17h23m, etc

UC05: Twice a day, first at 5h12 <==> Every 12h (period=12h), first at 5h12m

Specification
-------------

We can consider two basic types of scheduling:

1/ a one shot trigger at a predefined date
2/ a periodic trigger with a given duration

A real scheduling could be a composition of these two basic types.

We can configure a periodic scheduling by the parameters:

* period: duration of the cycle
* clock: the clock requirements
* round: the round requirements

UC01: period=24h, clock=5h12, round=nil (no rounding)
UC02: period=60m, clock=nil, round=0h10m
UC03: period=30m, clock=nil, round=0h10m
UC05: period=12h, clock=5h12, round=nil

The clock can be considered as a duration offset from the relative time zero (RTZ).
For a period greater that 24h (a day), the RTZ is 00h00m00s (begining of the day)

If the clock is given with 0 hour, then it should be considered as an offset until
the round hour given by the other parameters

*/

var logger = log.New(os.Stdout, "CRON:", log.LstdFlags)

func init() {
	logger.SetFlags(log.Ltime)
}

// LogActivate activates/deactivates the log output depending on the active
// parameter (true => activate, false => deactivate)
func LogActivate(active bool) {
	if active {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
}
