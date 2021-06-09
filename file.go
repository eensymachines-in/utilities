package utilities

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// FileAction : callback from the Watcher everytime file change is detected
type FileAction func(string) (interface{}, error)

// FileWatcher : sets up an infinite loop for an interval that will watch the file for chanages
// Upon finding a change - will callback where the client code can do what it needs to
// lower the tick time duration, faster is the file change detect, but higher is the load on CPU
// FileAction is called everytime the file is found to have been changed. FileAction is customizable on the client side
func FileWatcher(path string, cancel, errx chan interface{}, tick time.Duration, fa FileAction) (chan interface{}, func()) {
	out := make(chan interface{}, 2)
	return out, func() {
		defer close(out)
		for {
			initialStats, err := os.Stat(string(path))
			if err != nil {
				errx <- fmt.Errorf("Failed to read %s : Cannot start watcher loop", path)
				return
			}
			select {
			case <-cancel:
				// incase the client cancels we exit from here
				log.Warn("Now closing file watcher")
				return
			case <-time.After(tick):
				// check for the file stats if they've changed in anyway
				// else we just continue the loop
				stats, err := os.Stat(string(path))
				if err != nil {
					log.Error(err)
					errx <- fmt.Errorf("Failed to read new file stats %s: %s", path, err)
					return
				}
				if initialStats.Size() != stats.Size() || initialStats.ModTime() != stats.ModTime() {
					// this is when the file is found modified, time to stop schedule loop and restart again
					log.Info("File stats changed, now executing file action")
					result, err := fa(path)
					if err != nil {
						errx <- fmt.Errorf("Failed file action %s: %s", path, err)
						continue
					} else {
						out <- result
					}
				}
			}
		}
	}
}
