package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var logMu sync.Mutex
var dateFormat = "2006-01-02 15:04:05"

func InitLogger() {
	log.SetFlags(0)
	LogInfo("Logger initialized")
	LogInfo("OS: %s", runtime.GOOS)
	LogInfo("Architecture: %s", runtime.GOARCH)
}

func Log(logType string, format string, args ...interface{}) {
	logMu.Lock()
	defer logMu.Unlock()
	log.Printf(
		"[%s]  %s | %s",
		logType,
		time.Now().Format(dateFormat),
		fmt.Sprintf(format, args...),
	)
}

func LogInfo(format string, args ...interface{}) {
	Log("INFO", format, args...)
}

func LogWarn(format string, args ...interface{}) {
	Log("WARN", format, args...)
}

func LogError(format string, args ...interface{}) {
	Log("ERROR", format, args...)
}
