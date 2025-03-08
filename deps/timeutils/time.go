package timeutils

import "time"

var FoundingTime int64 = 1740230594
var FoundingTimeUTC = time.Unix(FoundingTime, 0).UTC()

type TimeNow func() time.Time

func (t TimeNow) Tick(d time.Duration) time.Time {
	return t().Add(d)
}
