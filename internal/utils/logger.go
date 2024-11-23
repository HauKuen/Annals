package utils

import (
	"time"

	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
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
	writer := &lumberjack.Logger{
		Filename: logDir + "/app.log",
		MaxSize:  500, // megabytes
		MaxAge:   31,  // days
		Compress: true,
	}

	// Set output to both custom writer and stdout
	multiWriter := io.MultiWriter(writer, os.Stdout)
	Log.SetOutput(multiWriter)

	Log.Info("Logger initialized. Daily log files will be created in ", logDir)
}

// LoggerMiddleware returns a Gin middleware for logging requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		Log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"duration": duration,
		}).Info("Request completed")
	}
}
