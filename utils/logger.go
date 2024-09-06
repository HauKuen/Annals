package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set log level based on the AppMode
	if AppMode == "debug" {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set up log directory
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Log.Fatal("Failed to create log directory:", err)
	}

	// Create a custom writer that rotates logs daily
	writer := &dailyRotateWriter{
		dir:        logDir,
		nameFormat: "2006-01-02.log",
	}

	// Set output to both custom writer and stdout
	multiWriter := io.MultiWriter(writer, os.Stdout)
	Log.SetOutput(multiWriter)

	// Start a goroutine to clean up old log files
	go cleanupOldLogs(logDir, 30)

	Log.Info("Logger initialized. Daily log files will be created in ", logDir)
}

type dailyRotateWriter struct {
	dir        string
	nameFormat string
	current    *lumberjack.Logger
}

func (w *dailyRotateWriter) Write(p []byte) (n int, err error) {
	if w.current == nil || w.current.Filename != w.getFilename() {
		if w.current != nil {
			w.current.Close()
		}
		w.current = &lumberjack.Logger{
			Filename: w.getFilename(),
			MaxSize:  500, // megabytes
			MaxAge:   31,  // days
			Compress: true,
		}
	}
	return w.current.Write(p)
}

func (w *dailyRotateWriter) getFilename() string {
	return filepath.Join(w.dir, time.Now().Format(w.nameFormat))
}

func cleanupOldLogs(dir string, maxAge int) {
	for {
		time.Sleep(24 * time.Hour) // Run once a day
		threshold := time.Now().AddDate(0, 0, -maxAge)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.ModTime().Before(threshold) {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Failed to remove old log file %s: %v\n", path, err)
				} else {
					fmt.Printf("Removed old log file: %s\n", path)
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking the path %q: %v\n", dir, err)
		}
	}
}
