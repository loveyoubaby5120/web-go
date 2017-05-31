package sentry

import (
	"fmt"
	"runtime"
	"strings"

	"jdy/pkg/base/trace"
	"jdy/pkg/util/errs"

	"github.com/getsentry/raven-go"
	"github.com/golang/glog"
)

var client *raven.Client
var errorLog = func() trace.T {
	t, _ := trace.NewEventLog("sentry", "error")
	return t
}()
var droppedLog = func() trace.T {
	t, _ := trace.NewEventLog("sentry", "dropped")
	return t
}()
var dataErrorLog = func() trace.T {
	t, _ := trace.NewEventLog("sentry", "data_error")
	return t
}()

// Init intializes the global client. Do not call it twice.
func Init(dsn string) error {
	if client != nil {
		panic("sentry.Init shouldn't be called twice!")
	}

	if dsn == "" {
		return nil
	}
	var err error
	client, err = raven.New(dsn)
	client.DropHandler = func(p *raven.Packet) {
		// droppedLog.Printf("%s", string(p.JSON()))
	}
	return err
}

func capture(depth int, err error, kind string, log trace.T) string {
	if err == nil {
		return ""
	}
	glog.ErrorDepth(2+depth, kind, " : ", err)
	log.Printf("%s %v", errs.CallerInfo(1+depth), err)

	if client == nil {
		return "<error not recorded>"
	}
	packet := raven.NewPacket(
		err.Error(), raven.NewException(err, raven.NewStacktrace(1+depth, 3, nil)))
	eventID, _ := client.Capture(packet, map[string]string{"kind": kind})
	return eventID
}

// ErrorDepth sends the error to sentry asynchronously.
func ErrorDepth(depth int, err error) string {
	return capture(depth, err, "error", errorLog)
}

// Error a shortcut of ErrorDepth(0, ...).
func Error(err error) string {
	return ErrorDepth(1, err)
}

// DataErrorDepth sends a data inconsistency error to sentry asynchronously.
func DataErrorDepth(depth int, err error) string {
	return capture(depth, err, "data_error", dataErrorLog)
}

// DataError is a shortcut of DataErrorDepth(0, ...).
func DataError(err error) string {
	return DataErrorDepth(1, err)
}

// WXBotDown sends a message to sentry with it's kind equals "wxbot_down"
func WXBotDown(err error) {
	capture(1, err, "wxbot_down", errorLog)
}

func callerInfo(depth int) string {
	_, file, line, ok := runtime.Caller(1 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
