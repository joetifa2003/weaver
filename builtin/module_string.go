package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/vm"
)

func registerStringModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("strings", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"concat": vm.NewNativeFunction("concat", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					var res string
					for _, arg := range args.Args {
						if arg.IsError() {
							return arg, false
						}
						res += arg.String()
					}

					return vm.NewString(res), true
				}),
				"split": vm.NewNativeFunction("split", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}

					sepArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return sepArg, false
					}

					str := strArg.GetString()
					sep := sepArg.GetString()

					parts := helpers.SliceMap(strings.Split(str, sep), func(s string) vm.Value {
						return vm.NewString(s)
					})

					return vm.NewArray(parts), true
				}),
				"lower": vm.NewNativeFunction("lower", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					return vm.NewString(strings.ToLower(strArg.GetString())), true
				}),
				"upper": vm.NewNativeFunction("upper", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					return vm.NewString(strings.ToUpper(strArg.GetString())), true
				}),
				"trim": vm.NewNativeFunction("trim", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}

					str := strArg.GetString()
					return vm.NewString(strings.TrimSpace(str)), true
				}),
				"contains": vm.NewNativeFunction("contains", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					substrArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return substrArg, false
					}
					return vm.NewBoolean(strings.Contains(strArg.GetString(), substrArg.GetString())), true
				}),
				"startsWith": vm.NewNativeFunction("startsWith", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					prefixArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return prefixArg, false
					}
					return vm.NewBoolean(strings.HasPrefix(strArg.GetString(), prefixArg.GetString())), true
				}),
				"endsWith": vm.NewNativeFunction("endsWith", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					suffixArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return suffixArg, false
					}
					return vm.NewBoolean(strings.HasSuffix(strArg.GetString(), suffixArg.GetString())), true
				}),
				"fmt": vm.NewNativeFunction("fmt", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					formatStrArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return formatStrArg, false
					}
					formatStr := formatStrArg.GetString()
					formatArgs := args.Args[1:]

					var result strings.Builder
					argIndex := 0

					for i := 0; i < len(formatStr); i++ {
						char := formatStr[i]

						switch char {
						case '\\':
							// Check for escape sequence
							if i+1 < len(formatStr) {
								nextChar := formatStr[i+1]
								if nextChar == '{' || nextChar == '}' || nextChar == '\\' {
									// Append the escaped character
									result.WriteByte(nextChar)
									i++ // Skip the next character
								} else {
									// Not a valid escape sequence, append the backslash literally
									result.WriteByte(char)
								}
							} else {
								// Backslash at the end of the string, append it literally
								result.WriteByte(char)
							}
						case '{':
							// Check if next char is '}' for placeholder
							if i+1 < len(formatStr) && formatStr[i+1] == '}' {
								if argIndex >= len(formatArgs) {
									// Not enough arguments provided for placeholders, append "{}" literally
									result.WriteString("{}")
								} else {
									arg := formatArgs[argIndex]
									if arg.IsError() { // Propagate errors from arguments
										return arg, false
									}
									result.WriteString(arg.String()) // Append string representation of the arg
									argIndex++
								}
								i++ // Skip the '}'
							} else {
								// It's just a literal '{', append it
								result.WriteByte(char)
							}
						default:
							// Append regular character
							result.WriteByte(char)
						}
					}

					return vm.NewString(result.String()), true
				}),
				"replace": vm.NewNativeFunction("replace", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					oldArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return oldArg, false
					}
					newArg, ok := args.Get(2, vm.ValueTypeString)
					if !ok {
						return newArg, false
					}

					count := -1 // Replace all by default
					if len(args.Args) > 3 {
						countArg, ok := args.Get(3, vm.ValueTypeNumber)
						if !ok {
							return countArg, false // Error if type is wrong
						}
						count = int(countArg.GetNumber())
					}

					return vm.NewString(strings.Replace(strArg.GetString(), oldArg.GetString(), newArg.GetString(), count)), true
				}),
				"substring": vm.NewNativeFunction("substring", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					startArg, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return startArg, false
					}

					str := strArg.GetString()
					start := int(startArg.GetNumber())
					end := len(str) // Default to end of string

					if len(args.Args) > 2 {
						endArg, ok := args.Get(2, vm.ValueTypeNumber)
						if !ok {
							return endArg, false // Error if type is wrong
						}
						end = int(endArg.GetNumber())
					}

					// Basic bounds checking
					if start < 0 {
						start = 0
					}
					if end > len(str) {
						end = len(str)
					}
					if start > end || start >= len(str) {
						return vm.NewString(""), true // Return empty string for invalid ranges
					}

					return vm.NewString(str[start:end]), true
				}),
				"indexOf": vm.NewNativeFunction("indexOf", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					substrArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return substrArg, false
					}

					return vm.NewNumber(float64(strings.Index(strArg.GetString(), substrArg.GetString()))), true
				}),
				"lastIndexOf": vm.NewNativeFunction("lastIndexOf", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					substrArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return substrArg, false
					}

					return vm.NewNumber(float64(strings.LastIndex(strArg.GetString(), substrArg.GetString()))), true
				}),
				"padStart": vm.NewNativeFunction("padStart", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					lengthArg, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return lengthArg, false
					}
					padStrArg := vm.NewString(" ") // Default pad string
					if len(args.Args) > 2 {
						padStrArg, ok = args.Get(2, vm.ValueTypeString)
						if !ok {
							return padStrArg, false
						}
					}

					str := strArg.GetString()
					targetLength := int(lengthArg.GetNumber())
					padStr := padStrArg.GetString()

					return vm.NewString(padString(str, targetLength, padStr, true)), true
				}),
				"padEnd": vm.NewNativeFunction("padEnd", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					strArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return strArg, false
					}
					lengthArg, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return lengthArg, false
					}
					padStrArg := vm.NewString(" ") // Default pad string
					if len(args.Args) > 2 {
						padStrArg, ok = args.Get(2, vm.ValueTypeString)
						if !ok {
							return padStrArg, false
						}
					}

					str := strArg.GetString()
					targetLength := int(lengthArg.GetNumber())
					padStr := padStrArg.GetString()

					return vm.NewString(padString(str, targetLength, padStr, false)), true
				}),
			},
		)
	})
}

// Helper function for padding
func padString(str string, targetLength int, padString string, padStart bool) string {
	if len(str) >= targetLength {
		return str
	}

	if padString == "" {
		padString = " " // Default pad string is space
	}

	padLen := targetLength - len(str)
	repeatCount := padLen / len(padString)
	remainingPad := padLen % len(padString)

	padding := strings.Repeat(padString, repeatCount) + padString[:remainingPad]

	if padStart {
		return padding + str
	}
	return str + padding
}
