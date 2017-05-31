package lu

import (
	"net/http"
	"reflect"
)

// BindAsHandler :: *Injector -> func(depdency) lu.Replyer -> result.
func BindAsHandler(inj *Injector, f interface{}) http.HandlerFunc {
	binded := inj.Bind(f)
	fT := reflect.TypeOf(f)

	var r Replyer
	if fT.NumOut() != 1 || fT.Out(0) != reflect.TypeOf(r) {
		panic("f must return a single value of type: Replyer.")
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		writer := binded(resp, req, func() {})[0]
		writer.Call([]reflect.Value{reflect.ValueOf(resp)})
	}
}

// Public :: *Injector -> func(depedency) lu.Replyer -> result.
func Public(inj *Injector, f interface{}) http.HandlerFunc {
	return BindAsHandler(inj, f)
}
