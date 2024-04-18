package clock

import "time"

type NowTimeFunc func() time.Time

var NowFunc NowTimeFunc

func ResetClockImplementation() {
	NowFunc = func() time.Time {
		return time.Now()
	}
}

func Now() time.Time {
	return NowFunc().UTC()
}
