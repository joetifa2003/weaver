package builtin

import (
	"encoding/json"

	"github.com/joetifa2003/weaver/vm"
)

func valufiyJSON(v interface{}) vm.Value {
	switch v := v.(type) {
	case string:
		return vm.NewString(v)
	case bool:
		return vm.NewBool(v)
	case int:
		return vm.NewInt(v)
	case int8:
		return vm.NewInt(int(v))
	case int16:
		return vm.NewInt(int(v))
	case int32:
		return vm.NewInt(int(v))
	case float64:
		return vm.NewFloat(v)
	case float32:
		return vm.NewFloat(float64(v))
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

func registerJSONModule(builder *RegistryBuilder) {
	builder.RegisterModule("json", map[string]vm.Value{
		"parse": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) vm.Value {
			if len(args) != 1 {
				panic("invalid number of arguments")
			}

			data := args[0].String()

			var result interface{}
			err := json.Unmarshal([]byte(data), &result)
			if err != nil {
				panic(err)
			}

			return valufiyJSON(result)
		}),
	})
}
