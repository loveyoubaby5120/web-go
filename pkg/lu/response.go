package lu

import (
	"fmt"
	"jdy/pkg/util/errs"
	"net/http"

	"github.com/golang/glog"

	"jdy/pkg/api"
)

func getCSRFToken(req *api.Request) string {
	return req.CSRFToken()
}

type Replyer func(res http.ResponseWriter)

// Error returns an error to user.
func Error(err error) Replyer {
	caller := errs.CallerInfo(1)
	return func(res http.ResponseWriter) {
		if e2, ok := err.(*errs.Error); ok {
			if e2.Kind == errs.Forbidden {
				glog.Infof("Forbid: %s %+v", caller, err)
				http.Error(res, "", http.StatusForbidden)
				return
			}
			if e2.Kind != errs.Internal {
				http.Error(res, e2.Error(), e2.Kind.HTTPStatusCode())
				return
			}
		}
		panic(fmt.Errorf("%s %v", caller, err))
	}
}
