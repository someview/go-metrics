//go:build !cgo || appengine
// +build !cgo appengine

package state

func numCgoCall() int64 {
	return 0
}
