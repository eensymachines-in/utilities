package utilities

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

// SetUpLog : here this is set up the most vanilla log preferences to get any module started
// Once upon the program exit return function can help close the log file if its not nil
// For most of thhe projects that we do here at eensymachines such logging setup in quite useful
func SetUpLog() func() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
		PadLevelText:  true,
	})
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel) // default is debug , the most verbose logging
	// file logs or logging to standard output
	val := os.Getenv("FLOG")
	var f *os.File
	var err error
	if val == "1" {
		f, err = os.OpenFile(os.Getenv("LOGF"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664) // file for logs
		if err != nil {
			log.SetOutput(os.Stdout) // error in opening log file
			log.Warn("Failed to open log file, log output set to stdout")
		}
		log.SetOutput(f) // log output set to file direction
		log.Infof("log output is set to file: %s", os.Getenv("LOGF"))

	} else {
		log.SetOutput(os.Stdout)
		log.Info("log output to stdout")
	}
	// how verbose would you want the logging to be
	val = os.Getenv("SILENT")
	if val == "1" {
		log.SetLevel(log.ErrorLevel) // for production
	} else {
		log.SetLevel(log.DebugLevel) // for development
	}
	return func() {
		sync.OnceFunc(func() {
			if f != nil {
				f.Close()
			}
		})
	}
}
