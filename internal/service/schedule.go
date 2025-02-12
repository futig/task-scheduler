package service

import (
	"time"

	"github.com/google/uuid"
)

func GetCurrentSchedule(chatID int64) (string, error) {
	return GetSchedule(chatID, time.Now().Local().Weekday())
}

func GetSchedule(chatID int64, weekday time.Weekday) (string, error) {

}

func GetCurrentTask(chatID int64) (string, error) {

}

func CreateSchedule(chatID int64, weekday time.Weekday) error {

}

func UpdateSchedule(chatID int64, weekday time.Weekday) error {

}

func DeleteSchedule(chatID int64, weekday time.Weekday) error {

}

func DeleteScheduleItem(chatID int64, weekday time.Weekday, id uuid.UUID) error {

}
