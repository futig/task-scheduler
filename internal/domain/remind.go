package domain

import "time"

type Remind struct {
	TimingId string
	TaskId   string
	Type     RemindType
	Time     time.Time
}
