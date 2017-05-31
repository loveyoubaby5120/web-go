package trace

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/trace"
)

// T provides uniform interface to track short/long-lived object.
type T interface {
	Printf(format string, a ...interface{})
	SetError()
}

type traceWrap struct {
	t trace.Trace
}

func (t traceWrap) Printf(format string, a ...interface{}) {
	t.t.LazyPrintf(format, a...)
}

func (t traceWrap) finish() {
	t.t.Finish()
}

func (t traceWrap) SetError() {
	t.t.SetError()
}

// New creates a new short-lived trace.
func New(family, title string) (T, func()) {
	t := traceWrap{trace.New(family, title)}
	t.t.SetMaxEvents(100)
	return t, t.t.Finish
}

type eventLogWrap struct {
	evl trace.EventLog
}

func (e eventLogWrap) Printf(format string, a ...interface{}) {
	e.evl.Printf(format, a...)
}

func (e eventLogWrap) SetError() {
	// TODO(guye): Log linenum heres.
	e.evl.Errorf("error state was set")
}

// NewEventLog creates a new long-lived trace.
func NewEventLog(family, title string) (T, func()) {
	t := eventLogWrap{trace.NewEventLog(family, title)}
	return t, t.evl.Finish
}

type noopTrace struct{}

func (t noopTrace) Printf(format string, a ...interface{}) {}
func (t noopTrace) SetError()                              {}

// Noop is a trace that does nothing.
var Noop T = &noopTrace{}

// TODO is a trace that needs to be replaced with a real one.
var TODO T = &noopTrace{}

type glogWrap struct{}

func (t glogWrap) Printf(format string, a ...interface{}) {
	glog.InfoDepth(2, fmt.Sprintf(format, a...))
}
func (t glogWrap) SetError() {}

// GLogTrace should only be used for unit test.
var GLogTrace T = &glogWrap{}
