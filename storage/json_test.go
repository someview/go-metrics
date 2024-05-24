package storage

import (
	"bytes"
	"encoding/json"
	"github.com/someview/go-metrics"
	"github.com/someview/go-metrics/counter"
	"testing"
)

func TestRegistryMarshallJSON(t *testing.T) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	r := metrics.NewRegistry()
	r.Register("counter", counter.NewCounter())
	enc.Encode(r)
	if s := b.String(); "{\"counter\":{\"count\":0}}\n" != s {
		t.Fatalf(s)
	}
}

func TestRegistryWriteJSONOnce(t *testing.T) {
	r := metrics.NewRegistry()
	r.Register("counter", counter.NewCounter())
	b := &bytes.Buffer{}
	WriteJSONOnce(r, b)
	if s := b.String(); s != "{\"counter\":{\"count\":0}}\n" {
		t.Fail()
	}
}
