package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogs() {
	logFile, err := os.OpenFile("./logs/logfile.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log := slog.New(slog.NewTextHandler(logFile, nil))

	Log = log
}
