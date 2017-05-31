package api

import (
	"net/http"
	"strings"
)

type Request struct {
	sessionKey string
	userAgent  string
	csrfToken  string
}

func NewRequestFromRequest(req *http.Request) *Request {
	return &Request{
		sessionKey: req.Header.Get("X-A2-REQUEST-SESSION-KEY"),
		userAgent:  strings.ToLower(req.Header.Get("X-A2-REQUEST-USER-AGENT")),
		csrfToken:  req.Header.Get("X-A2-REQUEST-CSRF-TOKEN"),
	}
}

func (r *Request) SessionKey() string {
	return r.sessionKey
}

func (r *Request) CSRFToken() string {
	return r.csrfToken
}

func (r *Request) FromAndroid() bool {
	return strings.Contains(r.userAgent, "android")
}

func (r *Request) FromIOS() bool {
	return strings.Contains(r.userAgent, "iphone") || strings.Contains(r.userAgent, "ipad")
}

func (r *Request) FromMobile() bool {
	return r.FromAndroid() || r.FromIOS()
}

func (r *Request) FromApplySquare() bool {
	return strings.Contains(r.userAgent, "applysquare")
}

func (r *Request) FromQQOrWeChat() bool {
	return strings.Contains(r.userAgent, "micromessenger") || strings.Contains(r.userAgent, "qq")
}

func (r *Request) FromKnownSpider() bool {
	spiders := []string{
		// baidu
		"baiduspider",
		// google
		"googlebot",
		// yahoo
		"yahoo! slurp",
		// 新浪爱问
		"iaskspider",
		// 搜狗
		"sogou web spider",
		"sogou inst spider",
		"sogou spider",
		// 网易有道
		"youdaobot",
		// msn
		"msnbot",
		// soso
		"sosospider",
		// bing
		"bingbot",
		// 360
		"360spider",
		// 阿里
		"yisouspider",
	}
	for _, spider := range spiders {
		if strings.Contains(r.userAgent, spider) {
			return true
		}
	}
	return false
}
