package builtin

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/joetifa2003/weaver/vm"
)

func registerHTTPModule(builder *vm.RegistryBuilder) {
	m := map[string]vm.Value{
		"request": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			// Get the options object
			optionsArg, ok := args.Get(0, vm.ValueTypeObject)
			if !ok {
				return optionsArg
			}

			options := optionsArg.GetObject()

			// Get required URL
			urlVal, ok := options["url"]
			if !ok {
				return vm.NewError("[http.request]: missing required field 'url'", vm.Value{})
			}
			url, ok := vm.CheckValueType(urlVal, vm.ValueTypeString)
			if !ok {
				return url
			}

			// Get required method
			methodVal, ok := options["method"]
			if !ok {
				return vm.NewError("[http.request]: missing required field 'method'", vm.Value{})
			}
			method, ok := vm.CheckValueType(methodVal, vm.ValueTypeString)
			if !ok {
				return method
			}

			// Create and make the request
			req, err, ok := createRequest(strings.ToUpper(method.GetString()), url.GetString(), vm.NativeFunctionArgs{optionsArg})
			if !ok {
				return err
			}

			return makeRequest(req)
		}),

		"newRouter": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			router := chi.NewRouter()
			return vm.NewNativeObject(router, makeRouterMethods(router))
		}),

		"listenAndServe": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			addrArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return addrArg
			}

			routerArg, ok := args.Get(1, vm.ValueTypeNativeObject)
			if !ok {
				return routerArg
			}

			router := routerArg.GetNativeObject().Obj.(*chi.Mux)

			err := http.ListenAndServe(addrArg.GetString(), router)
			if err != nil {
				return vm.NewErrFromErr(err)
			}

			return vm.Value{}
		}),
	}

	for _, method := range []string{"get", "post", "put", "delete"} {
		m[method] = vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			urlArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return urlArg
			}

			req, err, ok := createRequest(strings.ToUpper(method), urlArg.GetString(), args)
			if !ok {
				return err
			}

			return makeRequest(req)
		})
	}

	builder.RegisterModule("http", vm.NewObject(m))
}

func createRequest(method string, url string, args vm.NativeFunctionArgs) (*http.Request, vm.Value, bool) {
	var err error

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, vm.NewErrFromErr(err), false
	}

	if len(args) <= 1 {
		return req, vm.Value{}, true
	}

	optionsArg, ok := args.Get(1, vm.ValueTypeObject)
	if !ok {
		return nil, optionsArg, false
	}

	options := optionsArg.GetObject()

	if body, ok := options["body"]; ok {
		stringifiedBody, ok := stringify(body)
		if !ok {
			return nil, stringifiedBody, false
		}

		req.Body = io.NopCloser(strings.NewReader(stringifiedBody.GetString()))
	}

	headers, ok := options["headers"]
	if ok {
		headers, ok := vm.CheckValueType(headers, vm.ValueTypeObject)
		if !ok {
			return nil, headers, false
		}

		for key, val := range headers.GetObject() {
			req.Header.Add(key, val.String())
		}
	}

	return req, vm.Value{}, true
}

func makeRequest(req *http.Request) vm.Value {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return vm.NewError(err.Error(), vm.Value{})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return vm.NewError(err.Error(), vm.Value{})
	}

	return createResponseObject(resp, body)
}

func createResponseObject(resp *http.Response, body []byte) vm.Value {
	headers := make(map[string]vm.Value)
	for key, values := range resp.Header {
		headers[key] = vm.NewString(strings.Join(values, ", "))
	}

	obj := vm.NewObject(map[string]vm.Value{
		"statusCode": vm.NewNumber(float64(resp.StatusCode)),
		"status":     vm.NewString(resp.Status),
		"headers":    vm.NewObject(headers),
		"body":       vm.NewString(string(body)),
	})
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return obj
	}

	return vm.NewError("[http]: non 2xx response", obj)
}

func makeRouterMethods(router *chi.Mux) map[string]vm.Value {
	return map[string]vm.Value{
		"get": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			handlerArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return handlerArg
			}

			router.Get(pathArg.GetString(), func(w http.ResponseWriter, r *http.Request) {
				runHandler(v, handlerArg, r, w)
			})

			return vm.Value{}
		}),
		"post": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			handlerArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return handlerArg
			}

			router.Post(pathArg.GetString(), func(w http.ResponseWriter, r *http.Request) {
				runHandler(v, handlerArg, r, w)
			})

			return vm.Value{}
		}),
		"put": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			handlerArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return handlerArg
			}

			router.Put(pathArg.GetString(), func(w http.ResponseWriter, r *http.Request) {
				runHandler(v, handlerArg, r, w)
			})

			return vm.Value{}
		}),
		"delete": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			handlerArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return handlerArg
			}

			router.Delete(pathArg.GetString(), func(w http.ResponseWriter, r *http.Request) {
				runHandler(v, handlerArg, r, w)
			})

			return vm.Value{}
		}),
	}
}

type response struct {
	headers map[string]string
	status  int
}

func runHandler(v *vm.VM, handlerArg vm.Value, req *http.Request, w http.ResponseWriter) {
	response := response{
		headers: map[string]string{},
		status:  http.StatusOK,
	}

	task := runFunc(v, handlerArg, makeRequestObject(req), makeResponseObject(&response))
	val := task.Wait()
	if val.IsError() {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(val.String())
		return
	}

	for key, val := range response.headers {
		w.Header().Set(key, val)
	}

	switch val.VType {
	case vm.ValueTypeObject:
		w.Header().Set("Content-Type", "application/json")
		str, ok := stringify(val)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(response.status)
		w.Write([]byte(str.String()))
	default:
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/plain")
		}
		w.WriteHeader(response.status)
		w.Write([]byte(val.String()))
	}
}

func makeRequestObject(req *http.Request) vm.Value {
	return vm.NewObject(map[string]vm.Value{
		"method": vm.NewString(req.Method),
		"url":    vm.NewString(req.URL.String()),
		"path":   vm.NewString(req.URL.Path),
		"getQuery": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return keyArg
			}

			return vm.NewString(req.URL.Query().Get(keyArg.GetString()))
		}),
		"getFormValue": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return keyArg
			}

			return vm.NewString(req.FormValue(keyArg.GetString()))
		}),
		"getParam": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return keyArg
			}

			return vm.NewString(chi.URLParam(req, keyArg.GetString()))
		}),
	})
}

func makeResponseObject(resp *response) vm.Value {
	return vm.NewObject(map[string]vm.Value{
		"setHeader": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return keyArg
			}

			valArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return valArg
			}

			resp.headers[keyArg.GetString()] = valArg.GetString()

			return vm.Value{}
		}),
		"setStatus": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			statusArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return statusArg
			}

			resp.status = int(statusArg.GetNumber())

			return vm.Value{}
		}),
	})
}
