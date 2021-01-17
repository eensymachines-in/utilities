package utilities

import "sync"

// SafeCloseChn : channels need to be closed once, this can help to do that not more than once
// closure here deals with closing a channnel safely once using a niladic function
// when a channel has multiple points from where its closed,
// these multiple points cannot be predicted as in sequence
func SafeCloseChn(ch chan interface{}) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			close(ch)
		})
	}
}
