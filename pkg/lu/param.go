package lu

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"jdy/pkg/util/errs"
	"jdy/pkg/util/must"
)

var (
	paramDecoder = func() *schema.Decoder {
		d := schema.NewDecoder()
		d.IgnoreUnknownKeys(true)
		return d
	}()
)

type Param struct {
	req    *http.Request
	Vars   map[string]string
	errors []string
}

func NewParam(req *http.Request) *Param {
	p := &Param{req, mux.Vars(req), nil}

	var err error
	// Note: The django go proxy doesn't provide proper support for both file upload
	// and other form fields in form data, this is a bug but we don't have a use case
	// to fix it yet.
	if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		err = req.ParseMultipartForm(64 * 1024)
	} else {
		err = req.ParseForm()
	}
	if err != nil {
		p.AddError(fmt.Sprintf("failed to parse multipart form: %v", err))
	}
	return p
}

func (p *Param) Optional(key string, ret *string) *Param {
	v := p.req.Form[key]
	if len(v) > 0 {
		*ret = v[0]
	}
	return p
}

func (p *Param) OptionalInt(key string, ret *int) *Param {
	var s string
	p.Optional(key, &s)
	if p.Error() != nil {
		return p
	}
	if s == "" {
		return p
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		p.AddError(fmt.Sprintf("%s is not an integer: %s", key, err.Error()))
		return p
	}
	*ret = i
	return p
}

func (p *Param) AddError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Param) Required(key string, ret *string) *Param {
	v := p.req.Form[key]
	if len(v) == 0 || v[0] == "" {
		p.AddError(fmt.Sprintf("%s not set", key))
		return p
	}
	*ret = v[0]
	return p
}

func (p *Param) RequiredJSON(key string, j interface{}) *Param {
	str := ""
	p.Required(key, &str)
	if p.Error() != nil {
		return p
	}
	err := json.Unmarshal([]byte(str), j)
	if err != nil {
		p.AddError(fmt.Sprintf("%s: invalid json: %v", key, err))
		return p
	}
	return p
}

func (p *Param) Bool(key string, v *bool) *Param {
	str := ""
	p.Optional(key, &str)
	if p.Error() != nil {
		return p
	}
	if str == "" {
		return p
	}
	boo, err := strconv.ParseBool(str)
	if err != nil {
		p.AddError(fmt.Sprintf("%s: invalid bool: %v", key, err))
		return p
	}
	*v = boo
	return p
}

func (p *Param) Data(d interface{}) *Param {
	err := paramDecoder.Decode(d, p.req.Form)
	if err != nil {
		p.AddError(err.Error())
	}
	return p
}

func (p *Param) JSONBody(d interface{}) *Param {
	b, err := ioutil.ReadAll(p.req.Body)
	must.Must(err)
	err = json.Unmarshal(b, d)
	if err != nil {
		p.AddError(fmt.Sprintf("invalid body: %v", err.Error()))
	}
	return p
}

func (p *Param) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return errs.InvalidArgument.New(strings.Join(p.errors, "\n"))
}

func (p *Param) FormFile(key string, f *multipart.File, h *multipart.FileHeader) *Param {
	file, header, err := p.req.FormFile(key)
	if err != nil {
		p.AddError(fmt.Sprintf("%s: invalid bool: %v", key, err))
		return p
	}
	*f = file
	if h == nil {
		return p
	}
	*h = *header
	return p
}

func (p *Param) XMLBody(d interface{}) *Param {
	b, err := ioutil.ReadAll(p.req.Body)
	must.Must(err)
	err = xml.Unmarshal(b, d)
	if err != nil {
		p.AddError(fmt.Sprintf("invalid body: %v", err.Error()))
	}
	return p
}
