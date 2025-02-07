package domain

import "time"

type Timing struct {
	TimingId string
	TaskId string
	Type TimingType
	Time time.Time
}