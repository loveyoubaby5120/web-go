package i18n

import (
	"fmt"

	"jdy/pkg/api"
	"jdy/pkg/util/h"
)

func IntlURL(intl *Context, head string, tail string, params h.Params) string {
	path := fmt.Sprint("/", head, "-", intl.ShortLangCode(), "/")
	if tail != "" {
		path = path + tail + "/"
	}
	return h.URL(path, params)
}

func GetIntlRoutes(genRoutes func(*Context) []api.Route) []api.Route {
	return append(genRoutes(CnContext), genRoutes(EnContext)...)
}
