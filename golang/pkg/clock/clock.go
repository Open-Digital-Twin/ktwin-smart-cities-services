package clock

import "time"

type NowTimeFunc func() *time.Time

var NowFunc NowTimeFunc

func ResetClockImplementation() {
	NowFunc = func() *time.Time {
		now := time.Now()
		return &now
	}
}

func Now() *time.Time {
	now := NowFunc().UTC()
	return &now
}
