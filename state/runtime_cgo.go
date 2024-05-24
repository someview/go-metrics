//go:build cgo && !appengine
// +build cgo,!appengine

package state

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}
