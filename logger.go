package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type FileLogger struct {
	file   *os.File
	logger *log.Logger
}

func NewFileLogger(baseFilename string) *FileLogger {
	timestamp := time.Now().Format("02-01-2006 15-04-05")
	filename := fmt.Sprintf("%s_%s.log", baseFilename, timestamp)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &FileLogger{
		file:   file,
		logger: log.New(file, "", log.Ldate|log.Ltime),
	}
}

func (l *FileLogger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// -- Wails logger.Logger Interface Implementation --

func (l *FileLogger) Print(message string) {
	l.logger.Println(message)
}

func (l *FileLogger) Trace(message string) {
	l.logger.Println("TRC | " + message)
}

func (l *FileLogger) Debug(message string) {
	l.logger.Println("DBG | " + message)
}

func (l *FileLogger) Info(message string) {
	l.logger.Println("INF | " + message)
}

func (l *FileLogger) Warning(message string) {
	l.logger.Println("WRN | " + message)
}

func (l *FileLogger) Error(message string) {
	l.logger.Println("ERR | " + message)
}

func (l *FileLogger) Fatal(message string) {
	l.logger.Fatalln("FTL | " + message)
}
