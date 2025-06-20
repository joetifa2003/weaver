package builtin

import (
	"time"

	"github.com/joetifa2003/weaver/vm"
)

func registerTimeModule(builder *vm.RegistryBuilder) {

	builder.RegisterModule("time", func() vm.Value {
		m := map[string]vm.Value{
			// --- Time Constants ---
			"ansic":       vm.NewString(time.ANSIC),
			"unixDate":    vm.NewString(time.UnixDate),
			"rubyDate":    vm.NewString(time.RubyDate),
			"rfc822":      vm.NewString(time.RFC822),
			"rfc822z":     vm.NewString(time.RFC822Z),
			"rfc850":      vm.NewString(time.RFC850),
			"rfc1123":     vm.NewString(time.RFC1123),
			"rfc1123z":    vm.NewString(time.RFC1123Z),
			"rfc3339":     vm.NewString(time.RFC3339),
			"rfc3339Nano": vm.NewString(time.RFC3339Nano),
			"kitchen":     vm.NewString(time.Kitchen),
			"stamp":       vm.NewString(time.Stamp),
			"stampMilli":  vm.NewString(time.StampMilli),
			"stampMicro":  vm.NewString(time.StampMicro),
			"stampNano":   vm.NewString(time.StampNano),
			"dateTime":    vm.NewString(time.DateTime),
			"dateOnly":    vm.NewString(time.DateOnly),
			"timeOnly":    vm.NewString(time.TimeOnly),

			// --- Duration Constants ---
			"nanosecond":  vm.NewNumber(float64(time.Nanosecond)),
			"microsecond": vm.NewNumber(float64(time.Microsecond)),
			"millisecond": vm.NewNumber(float64(time.Millisecond)),
			"second":      vm.NewNumber(float64(time.Second)),
			"minute":      vm.NewNumber(float64(time.Minute)),
			"hour":        vm.NewNumber(float64(time.Hour)),

			// --- Time Functions ---
			"now": vm.NewNativeFunction("now", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				return vm.NewTime(time.Now()), true
			}),

			"unix": vm.NewNativeFunction("unix", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				secArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return secArg, false
				}
				nsecArg := vm.NewNumber(0) // Default nsec to 0
				if len(args.Args) > 1 {
					nsecArg, ok = args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return nsecArg, false
					}
				}
				return vm.NewTime(time.Unix(int64(secArg.GetNumber()), int64(nsecArg.GetNumber()))), true
			}),

			"unixMilli": vm.NewNativeFunction("unixMilli", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				milliArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return milliArg, false
				}
				return vm.NewTime(time.UnixMilli(int64(milliArg.GetNumber()))), true
			}),

			"unixMicro": vm.NewNativeFunction("unixMicro", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				microArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return microArg, false
				}
				return vm.NewTime(time.UnixMicro(int64(microArg.GetNumber()))), true
			}),

			"parse": vm.NewNativeFunction("parse", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				layoutArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return layoutArg, false
				}
				valueArg, ok := args.Get(1, vm.ValueTypeString)
				if !ok {
					return valueArg, false
				}

				t, err := time.Parse(layoutArg.GetString(), valueArg.GetString())
				if err != nil {
					return vm.NewErrFromErr(err), false
				}
				return vm.NewTime(t), true
			}),

			"since": vm.NewNativeFunction("since", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(time.Since(timeArg.GetTime()))), true
			}),

			"until": vm.NewNativeFunction("until", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(time.Until(timeArg.GetTime()))), true
			}),

			// --- Duration Functions ---
			"parseDuration": vm.NewNativeFunction("parseDuration", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationStrArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return durationStrArg, false
				}
				d, err := time.ParseDuration(durationStrArg.GetString())
				if err != nil {
					return vm.NewErrFromErr(err), false
				}
				return vm.NewNumber(float64(d)), true
			}),

			// --- Time Methods (as functions taking time as first arg) ---
			"add": vm.NewNativeFunction("add", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				durationArg, ok := args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				t := timeArg.GetTime()
				d := time.Duration(durationArg.GetNumber())
				return vm.NewTime(t.Add(d)), true
			}),

			"sub": vm.NewNativeFunction("sub", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				time1Arg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return time1Arg, false
				}
				time2Arg, ok := args.Get(1, vm.ValueTypeTime)
				if !ok {
					return time2Arg, false
				}
				t1 := time1Arg.GetTime()
				t2 := time2Arg.GetTime()
				return vm.NewNumber(float64(t1.Sub(t2))), true
			}),

			"addDate": vm.NewNativeFunction("addDate", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				yearsArg, ok := args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return yearsArg, false
				}
				monthsArg, ok := args.Get(2, vm.ValueTypeNumber)
				if !ok {
					return monthsArg, false
				}
				daysArg, ok := args.Get(3, vm.ValueTypeNumber)
				if !ok {
					return daysArg, false
				}
				t := timeArg.GetTime()
				return vm.NewTime(t.AddDate(int(yearsArg.GetNumber()), int(monthsArg.GetNumber()), int(daysArg.GetNumber()))), true
			}),

			"after": vm.NewNativeFunction("after", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				time1Arg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return time1Arg, false
				}
				time2Arg, ok := args.Get(1, vm.ValueTypeTime)
				if !ok {
					return time2Arg, false
				}
				t1 := time1Arg.GetTime()
				t2 := time2Arg.GetTime()
				return vm.NewBool(t1.After(t2)), true
			}),

			"before": vm.NewNativeFunction("before", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				time1Arg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return time1Arg, false
				}
				time2Arg, ok := args.Get(1, vm.ValueTypeTime)
				if !ok {
					return time2Arg, false
				}
				t1 := time1Arg.GetTime()
				t2 := time2Arg.GetTime()
				return vm.NewBool(t1.Before(t2)), true
			}),

			"equal": vm.NewNativeFunction("equal", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				time1Arg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return time1Arg, false
				}
				time2Arg, ok := args.Get(1, vm.ValueTypeTime)
				if !ok {
					return time2Arg, false
				}
				t1 := time1Arg.GetTime()
				t2 := time2Arg.GetTime()
				return vm.NewBool(t1.Equal(t2)), true
			}),

			"format": vm.NewNativeFunction("format", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				layoutArg, ok := args.Get(1, vm.ValueTypeString)
				if !ok {
					return layoutArg, false
				}
				t := timeArg.GetTime()
				return vm.NewString(t.Format(layoutArg.GetString())), true
			}),

			"isZero": vm.NewNativeFunction("isZero", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewBool(timeArg.GetTime().IsZero()), true
			}),

			"date": vm.NewNativeFunction("date", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				t := timeArg.GetTime()
				year, month, day := t.Date()
				return vm.NewArray([]vm.Value{
					vm.NewNumber(float64(year)),
					vm.NewNumber(float64(month)),
					vm.NewNumber(float64(day)),
				}), true
			}),

			"getYear": vm.NewNativeFunction("getYear", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Year())), true
			}),

			"getMonth": vm.NewNativeFunction("getMonth", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Month())), true
			}),

			"getDay": vm.NewNativeFunction("getDay", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Day())), true
			}),

			"weekday": vm.NewNativeFunction("weekday", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Weekday())), true
			}),

			"clock": vm.NewNativeFunction("clock", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				t := timeArg.GetTime()
				hour, min, sec := t.Clock()

				return vm.NewObject(map[string]vm.Value{
					"hour":   vm.NewNumber(float64(hour)),
					"minute": vm.NewNumber(float64(min)),
					"second": vm.NewNumber(float64(sec)),
				}), true
			}),

			"getHour": vm.NewNativeFunction("getHour", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Hour())), true
			}),

			"getMinute": vm.NewNativeFunction("getMinute", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Minute())), true
			}),

			"getSecond": vm.NewNativeFunction("getSecond", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Second())), true
			}),

			"getNanosecond": vm.NewNativeFunction("getNanosecond", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Nanosecond())), true
			}),

			"getUnixTime": vm.NewNativeFunction("getUnixTime", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().Unix())), true
			}),

			"getUnixMilliTime": vm.NewNativeFunction("getUnixMilliTime", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().UnixMilli())), true
			}),

			"getUnixMicroTime": vm.NewNativeFunction("getUnixMicroTime", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().UnixMicro())), true
			}),

			"getUnixNanoTime": vm.NewNativeFunction("getUnixNanoTime", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewNumber(float64(timeArg.GetTime().UnixNano())), true
			}),

			"utc": vm.NewNativeFunction("utc", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewTime(timeArg.GetTime().UTC()), true
			}),

			"local": vm.NewNativeFunction("local", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				return vm.NewTime(timeArg.GetTime().Local()), true
			}),

			// --- Duration Methods (as functions taking duration number as first arg) ---
			"getHours": vm.NewNativeFunction("getHours", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(d.Hours()), true
			}),

			"getMinutes": vm.NewNativeFunction("getMinutes", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(d.Minutes()), true
			}),

			"getSeconds": vm.NewNativeFunction("getSeconds", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(d.Seconds()), true
			}),

			"getMilliseconds": vm.NewNativeFunction("getMilliseconds", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(float64(d.Milliseconds())), true
			}),

			"getMicroseconds": vm.NewNativeFunction("getMicroseconds", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(float64(d.Microseconds())), true
			}),

			"getNanoseconds": vm.NewNativeFunction("getNanoseconds", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewNumber(float64(d.Nanoseconds())), true
			}),

			"getDurationString": vm.NewNativeFunction("getDurationString", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				durationArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return durationArg, false
				}
				d := time.Duration(durationArg.GetNumber())
				return vm.NewString(d.String()), true
			}),

			// --- Timezone Functions ---
			"parseInLocation": vm.NewNativeFunction("parseInLocation", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				layoutArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return layoutArg, false
				}
				valueArg, ok := args.Get(1, vm.ValueTypeString)
				if !ok {
					return valueArg, false
				}
				locNameArg, ok := args.Get(2, vm.ValueTypeString)
				if !ok {
					return locNameArg, false
				}

				loc, err := time.LoadLocation(locNameArg.GetString())
				if err != nil {
					return vm.NewErrFromErr(err), false
				}

				t, err := time.ParseInLocation(layoutArg.GetString(), valueArg.GetString(), loc)
				if err != nil {
					return vm.NewErrFromErr(err), false
				}
				return vm.NewTime(t), true
			}),

			"inLocation": vm.NewNativeFunction("inLocation", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				locNameArg, ok := args.Get(1, vm.ValueTypeString)
				if !ok {
					return locNameArg, false
				}

				loc, err := time.LoadLocation(locNameArg.GetString())
				if err != nil {
					return vm.NewErrFromErr(err), false
				}

				t := timeArg.GetTime()
				return vm.NewTime(t.In(loc)), true
			}),

			"getZone": vm.NewNativeFunction("getZone", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				timeArg, ok := args.Get(0, vm.ValueTypeTime)
				if !ok {
					return timeArg, false
				}
				t := timeArg.GetTime()
				name, offset := t.Zone()
				return vm.NewObject(map[string]vm.Value{
					"name":   vm.NewString(name),
					"offset": vm.NewNumber(float64(offset)), // Offset in seconds
				}), true
			}),
		}

		return vm.NewObject(m)
	})
}
