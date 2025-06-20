package builtin

import (
	"encoding/json"

	"github.com/joetifa2003/weaver/vm"
)

func registerJSONModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("json", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"parse": vm.NewNativeFunction("parse", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					dataArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return dataArg, false
					}

					data := dataArg.String()
					var result interface{}
					err := json.Unmarshal([]byte(data), &result)
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}

					return valufiyJSON(result), true
				}),
				"stringify": vm.NewNativeFunction("stringify", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					dataArg, ok := args.Get(0)
					if !ok {
						return dataArg, false
					}

					res, ok := stringify(dataArg)
					if !ok {
						return res, false
					}

					return res, true
				}),
			},
		)
	})
}

func stringify(v vm.Value) (vm.Value, bool) {
	val := goifyValue(v)
	b, err := json.Marshal(val)
	if err != nil {
		return vm.NewError(err.Error(), vm.Value{}), false
	}

	return vm.NewString(string(b)), true
}

func valufiyJSON(v interface{}) vm.Value {
	switch v := v.(type) {
	case string:
		return vm.NewString(v)
	case bool:
		return vm.NewBool(v)
	case int:
		return vm.NewNumber(float64(v))
	case int8:
		return vm.NewNumber(float64(v))
	case int16:
		return vm.NewNumber(float64(v))
	case int32:
		return vm.NewNumber(float64(v))
	case float64:
		return vm.NewNumber(v)
	case float32:
		return vm.NewNumber(float64(v))
	case map[string]interface{}:
		m := make(map[string]vm.Value)
		for k, v := range v {
			m[k] = valufiyJSON(v)
		}
		return vm.NewObject(m)
	case []interface{}:
		a := make([]vm.Value, len(v))
		for i, v := range v {
			a[i] = valufiyJSON(v)
		}
		return vm.NewArray(a)
	default:
		panic("valuify json: unsupported type")
	}
}

func goifyValue(v vm.Value) interface{} {
	switch v.VType {
	case vm.ValueTypeString:
		return v.GetString()

	case vm.ValueTypeBool:
		return v.GetBool()

	case vm.ValueTypeNumber:
		return v.GetNumber()

	case vm.ValueTypeObject:
		m := make(map[string]interface{})
		for k, v := range v.GetObject() {
			m[k] = goifyValue(v)
		}
		return m

	case vm.ValueTypeArray:
		a := make([]interface{}, len(*v.GetArray()))
		for i, v := range *v.GetArray() {
			a[i] = goifyValue(v)
		}

		return a

	default:
		return nil
	}
}
