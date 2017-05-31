package trace

import (
	"fmt"

	"github.com/golang/glog"
)

// Infof log and print to the trace.
func Infof(tr T, format string, a ...interface{}) {
	tr.Printf(format, a...)
	glog.InfoDepth(2, fmt.Sprintf(format, a...))
}

// Errorf log and print to the trace.
func Errorf(tr T, format string, a ...interface{}) {
	tr.Printf(format, a...)
	glog.ErrorDepth(2, fmt.Sprintf(format, a...))
}

// Warningf log and print to the trace.
func Warningf(tr T, format string, a ...interface{}) {
	tr.Printf(format, a...)
	glog.WarningDepth(2, fmt.Sprintf(format, a...))
}
