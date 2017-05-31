package lu

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"jdy/pkg/api"
	"jdy/pkg/i18n"
)

type injVal struct {
	constV      reflect.Value
	deferred    interface{}
	directParam int
}

type Injector struct {
	typeValueMap map[reflect.Type]injVal
	parent       *Injector
}

func InterfaceOf(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("Called inject.InterfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
	}

	return t
}

func NewInjectorWithOptions(enableIntl bool) *Injector {
	inj := &Injector{
		typeValueMap: make(map[reflect.Type]injVal),
	}
	var res http.ResponseWriter
	inj.Set(InterfaceOf(&res), injVal{directParam: 0})
	var req *http.Request
	inj.Set(reflect.TypeOf(req), injVal{directParam: 1})
	var next Next
	inj.Set(reflect.TypeOf(next), injVal{directParam: 2})
	inj.SetDeferred(NewParam)

	// A2 API specific.
	inj.SetDeferred(api.NewUserFromRequest)
	inj.SetDeferred(api.NewRequestFromRequest)
	if enableIntl {
		inj.SetDeferred(func(req *http.Request) *i18n.Context {
			return i18n.NewContext(req.Header.Get("X-A2-LANGUAGE"))
		})
	}
	inj.SetDeferred(NewTracer)
	return inj
}

func NewInjector() *Injector {
	return NewInjectorWithOptions(true)
}

func NewChildInjector(parent *Injector) *Injector {
	inj := &Injector{
		typeValueMap: make(map[reflect.Type]injVal),
	}
	inj.parent = parent
	return inj
}

type Next func()

// SetDeferred registers a given type that depends on the request.
// f function of the form:
//   func(...) T.
func (inj *Injector) SetDeferred(f interface{}) {
	funT := reflect.TypeOf(f)
	t := funT.Out(0)
	inj.Set(t, injVal{deferred: f})
}

func (inj *Injector) SetConst(v interface{}) {
	inj.Set(reflect.TypeOf(v), injVal{constV: reflect.ValueOf(v)})
}

func (inj *Injector) SetConstInterface(iptr interface{}, v interface{}) {
	inj.Set(InterfaceOf(iptr), injVal{constV: reflect.ValueOf(v)})
}

func (inj *Injector) Set(t reflect.Type, v injVal) {
	_, exist := inj.typeValueMap[t]
	if exist {
		panic(fmt.Sprintf("Value for type: %v is already set.", t))
	}
	inj.typeValueMap[t] = v
}

func (inj *Injector) getInjValue(typ reflect.Type) (injVal, bool) {
	v, exist := inj.typeValueMap[typ]
	if exist {
		return v, true
	}
	if inj.parent == nil {
		return injVal{}, false
	}
	return inj.parent.getInjValue(typ)
}

func (inj *Injector) RegisteredTypeNames() []string {
	var result []string
	for ty := range inj.typeValueMap {
		result = append(result, fmt.Sprint(ty))
	}
	if inj.parent != nil {
		result = append(result, inj.parent.RegisteredTypeNames()...)
	}
	return result
}

func (inj *Injector) Bind(f interface{}) func(resp http.ResponseWriter, req *http.Request, next Next) []reflect.Value {
	funV := reflect.ValueOf(f)
	funT := reflect.TypeOf(f)
	// NumIn the number of input parameter
	type paramInfo struct {
		constV         reflect.Value
		hasConstV      bool
		deferredDirect int
		deferredFunc   func(resp http.ResponseWriter, req *http.Request) reflect.Value
	}

	paramInfos := make([]paramInfo, funT.NumIn())
	for i := 0; i < len(paramInfos); i++ {
		argType := funT.In(i)
		v, exist := inj.getInjValue(argType)
		if !exist {
			panic(fmt.Sprintf("Value for type: %v is not set, set types: %s", argType, strings.Join(inj.RegisteredTypeNames(), ",")))
		}
		if v.constV.IsValid() {
			paramInfos[i].constV = v.constV
			paramInfos[i].hasConstV = true
			continue
		}
		if v.deferred != nil {
			binded := inj.Bind(v.deferred)
			paramInfos[i].deferredFunc = func(resp http.ResponseWriter, req *http.Request) reflect.Value {
				return binded(resp, req, func() { panic("shouldn't used") })[0]
			}
			continue
		}
		paramInfos[i].deferredDirect = v.directParam
	}

	return func(resp http.ResponseWriter, req *http.Request, next Next) []reflect.Value {
		params := make([]reflect.Value, funT.NumIn())
		for i := range paramInfos {
			if paramInfos[i].hasConstV {
				params[i] = paramInfos[i].constV
				continue
			}
			if f := paramInfos[i].deferredFunc; f != nil {
				params[i] = f(resp, req)
				continue
			}
			switch paramInfos[i].deferredDirect {
			case 0:
				params[i] = reflect.ValueOf(resp)
			case 1:
				params[i] = reflect.ValueOf(req)
			case 2:
				params[i] = reflect.ValueOf(next)
			}
		}
		return funV.Call(params)
	}
}
