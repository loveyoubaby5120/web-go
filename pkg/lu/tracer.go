package lu

import (
	"net/http"

	"jdy/pkg/api"
	"jdy/pkg/base/trace"
)

// Tracer is a helper to create trace for http request.
type Tracer struct {
	req *http.Request
}

// NewTracer creates a new tracer for request.
func NewTracer(req *http.Request) *Tracer {
	return &Tracer{req}
}

// TraceWithTitle creates a trace with given family and title.
func (h *Tracer) TraceWithTitle(family, title string) (trace.T, func()) {
	// We can do some sampling here.
	tr, done := trace.New(family, title)
	tr.Printf("URL: %s", h.req.URL.String())
	return tr, func() {
		tr.Printf("Finish.")
		done()
	}
}

// Trace creates a trace with its own family.
func (h *Tracer) Trace(family string) (trace.T, func()) {
	if api.IsDev(h.req) && h.req.URL.Query().Get("manual_trace") != "" {
		return h.TraceWithTitle("manual", family)
	}
	return h.TraceWithTitle(family, "run")
}

// TraceMisc creates a trace within the "misc" family.
func (h *Tracer) TraceMisc(title string) (trace.T, func()) {
	return h.TraceWithTitle("Misc", title)
}

// TraceBulky creates a trace within the "bulky" family.
func (h *Tracer) TraceBulky(title string) (trace.T, func()) {
	return h.TraceWithTitle("Bulky", title)

}
