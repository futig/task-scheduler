package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/storage"
	"github.com/futig/task-scheduler/internal/storage/temp_storage"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var variablesInited bool

type AppConfig struct {
	MinWorkers           int
	MaxWorkers           int
	WorkersCheckInterval time.Duration
	BusyThreshold        int
	UpdatesQueueSize     int
	RemindsQueueSize     int
	ScaleUpThreshold     int
	ScaleDownThreshold   int
	RemindsCheckInterval time.Duration
}

func NewAppConfig() AppConfig {
	initEnvVariables()
	return AppConfig{
		MinWorkers:           getIntVar("MIN_WORKERS"),
		MaxWorkers:           getIntVar("MAX_WORKERS"),
		BusyThreshold:        getIntVar("BUSY_THRESHOLD"),
		RemindsQueueSize:     getIntVar("REMINDS_QUEUE_SIZE"),
		ScaleUpThreshold:     getIntVar("SCALE_UP_THRESHOLD"),
		ScaleDownThreshold:   getIntVar("SCALE_DOWN_THRESHOLD"),
		RemindsCheckInterval: time.Duration(getIntVar("CHECK_REMINDS_INTERVAL")) * time.Minute,
		WorkersCheckInterval: time.Duration(getIntVar("WORKERS_CHECK_INTERVAL")) * time.Second,
	}
}

type WorkflowConfig struct {
	RemindsCh    chan domain.TaskRemind
	UpdatesCh    tgbotapi.UpdatesChannel
	StopWorkerCh chan struct{}
	Mu           sync.Mutex
	Storage      storage.Storage
}

func NewWorkflowConfig() WorkflowConfig {
	initEnvVariables()
	return WorkflowConfig{
		RemindsCh:    make(chan domain.TaskRemind, getIntVar("REMINDS_QUEUE_SIZE")),
		StopWorkerCh: make(chan struct{}),
		Storage:      &tempstorage.TempStorageContext{},
	}
}

func (w *WorkflowConfig) CloseCh() {
	close(w.RemindsCh)
	close(w.StopWorkerCh)
}

func initEnvVariables() {
	if !variablesInited {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(".env file not found or can't be loaded")
		}
		variablesInited = true
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
