package enums

type UserState int

const (
	MainMenu RemindType = iota
	ChoosingDOWToGet
	ChoosingDOWToUpdate
	UpdatingDaySchedule
	UpdatingItem
	UpdatingItemDescription
	UpdatingItemTime
	UpdatingItemReminds
)
