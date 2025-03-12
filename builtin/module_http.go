package builtin

import (
	"io"
	"net/http"
	"strings"

	"github.com/joetifa2003/weaver/vm"
)

func registerHTTPModule(builder *RegistryBuilder) {
	m := map[string]vm.Value{}
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
	builder.RegisterModule("http", m)
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

	return vm.NewObject(map[string]vm.Value{
		"status":  vm.NewNumber(float64(resp.StatusCode)),
		"headers": vm.NewObject(headers),
		"body":    vm.NewString(string(body)),
	})
}
