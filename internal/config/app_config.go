package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	MinWorkers         int
	MaxWorkers         int
	BusyThreshold      int
	QueueSize          int
	ScaleUpThreshold   int
	ScaleDownThreshold int
}

func NewAppConfig() AppConfig{
	return AppConfig{
		MinWorkers: getIntVar("MIN_WORKERS"),
		MaxWorkers: getIntVar("MAX_WORKERS"),
		BusyThreshold: getIntVar("BUSY_THRESHOLD"),
		QueueSize: getIntVar("QUEUE_SIZE"),
		ScaleUpThreshold: getIntVar("SCALE_UP_THRESHOLD"),
		ScaleDownThreshold: getIntVar("SCALE_DOWN_THRESHOLD"),
	}
}

func getIntVar(key string) int {
	rawVal := os.Getenv(key)
	if rawVal == "" {
		log.Fatal(fmt.Sprintf("variable %s is not specified", key))
	}

	val, err := strconv.Atoi(rawVal)
	if err != nil {
		log.Fatal(fmt.Sprintf("variable %s must be integer", key))
	}

	return val
}
