package api

import (
	"encoding/base64"
	"fmt"
	"jdy/pkg/base/sentry"
	"jdy/pkg/common"
	"net/http"
)

type User struct {
	id          string
	permissions map[string]bool
	name        string
	email       string
	authSite    string
}

func NewUserFromRequest(req *http.Request) *User {
	id := req.Header.Get("X-A2-USER-UUID")
	authSite := req.Header.Get("X-A2-USER-AUTH-SITE")
	bname := req.Header.Get("X-A2-USER-NAME")
	name := ""
	if bname != "" {
		namebyte, err := base64.StdEncoding.DecodeString(bname)
		if err != nil {
			sentry.Error(fmt.Errorf("unable to decode user name %s", bname))
		}
		name = string(namebyte)
	}

	email := req.Header.Get("X-A2-USER-EMAIL")
	permissions := make(map[string]bool)

	for _, p := range common.SplitAndTrim(req.Header.Get("X-A2-USER-PERMISSIONS"), ",") {
		permissions[p] = true
	}
	return &User{
		id:          id,
		permissions: permissions,
		name:        name,
		email:       email,
		authSite:    authSite,
	}
}

func IsDev(r *http.Request) bool {
	for _, p := range common.SplitAndTrim(r.Header.Get("X-A2-USER-PERMISSIONS"), ",") {
		if p == "dev" {
			return true
		}
	}
	return false
}
