package utilities

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// SetUpLog : here this is set up the most vanilla log preferences to get any module started
func SetUpLog() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // default is info level, if verbose then trace
}

// CustomLog : depending on the flags sent in this can alter the log levels and redirect the output
// send in the logfile path so that incase of file logging this can create and and setupout to the file
func CustomLog(flog, fverbose bool, logFile string) {
	if flog {
		lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Error(err)
			log.Error("Failed to connect to log file, kindly check the privileges")
			// If this fails to connect to log file, it defaults to sending the log output to stdout
		} else {
			log.Infof("Check log file %s for entries", logFile)
			log.SetOutput(lf)
			defer lf.Close()
		}
	}
	if fverbose {
		log.Info("Verbose logging")
		log.SetLevel(log.TraceLevel)
	}
}
