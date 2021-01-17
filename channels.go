package utilities

import "sync"

// SafeCloseChn : channels need to be closed once, this can help to do that not more than once
func SafeCloseChn(ch chan interface{}) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			close(ch)
		})
	}
}
