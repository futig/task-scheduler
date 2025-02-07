package domain

type TimingType int

const (
	Reminder TimingType = iota
	Event
)