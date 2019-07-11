package mrlog

import "time"

type Clock struct{}

func (_ *Clock) Now() time.Time {
	return time.Now()
}
