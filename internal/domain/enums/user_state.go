package enums

type UserState int

const (
	MainMenu UserState = iota
	
	ChooseWeekdayGet
	ChooseWeekdayUpdate

	UpdatingDaySchedule

	UpdateTask
	UpdateTaskDescription
	UpdateTaskTime
	UpdateTaskReminds

	ShowCurrentSchedule
	ShowCurrentTask
	ShowOtherDaySchedule

	DeleteSchedule
	DeleteTask

	TaskNotFound
)
