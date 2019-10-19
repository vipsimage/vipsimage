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
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/vipsimage/vipsimage/route"
	"github.com/vipsimage/vipsimage/rule"
)

func init() {
	// init config
	logPath := `/data/logs`
	if !gin.IsDebugging() {
		viper.AddConfigPath("/data")
	} else {
		logPath = `data/logs`
		viper.AddConfigPath("data")
	}
	viper.SetConfigName("vipsimage")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Check log dir
	logDir := path.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Set default logrus time format
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var runtimeLog io.Writer = &lumberjack.Logger{
		Filename:   logPath + "/runtime.log",
		MaxSize:    64, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
	}

	var accessLog io.Writer = &lumberjack.Logger{
		Filename:   logPath + "/access.log",
		MaxSize:    64, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
	}

	// Debug write to stdout
	if gin.IsDebugging() {
		runtimeLog = os.Stdout
		accessLog = os.Stdout
	}

	logrus.SetOutput(runtimeLog)
	gin.DefaultWriter = accessLog
}

func main() {
	// Parse vipsimage.toml operation-rule
	rule.Init()

	// Bind host port
	addr := viper.GetString("vipsimage.bind")
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
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
