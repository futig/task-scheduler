package time

import "time"

func CurrentTimeToMinutes() int {
	curTime := time.Now().Local()
	totalMinutes := curTime.Hour()*60 + curTime.Minute()
	return totalMinutes
}
