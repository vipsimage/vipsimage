package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/vipsimage/vipsimage/route"
	"github.com/vipsimage/vipsimage/utils/conf"
)

func init() {
	logPath := `/app/data/logs`
	if gin.IsDebugging() {
		logPath = "data/logs"
	}
	logDir := path.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var runtimeLog io.Writer = &lumberjack.Logger{
		Filename:   conf.Getenv("RUNTIME_PATH", "data/logs/car-api.log"),
		MaxSize:    64, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
	}

	var accessLog io.Writer = &lumberjack.Logger{
		Filename:   conf.Getenv("ACCESS_PATH", "data/logs/access.log"),
		MaxSize:    64, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
	}

	if gin.IsDebugging() {
		runtimeLog = os.Stdout
		accessLog = os.Stdout
	}

	logrus.SetOutput(runtimeLog)
	gin.DefaultWriter = accessLog
}

func main() {
	addr := conf.Getenv("SERVER_ADDRESS", "0.0.0.0:80")

	fmt.Println("bind addr: ", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: route.Route(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
