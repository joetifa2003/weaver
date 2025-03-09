package builtin

import "github.com/joetifa2003/weaver/vm"

type RegistryBuilder struct {
	modules map[string]vm.Value
	funcs   map[string]vm.Value
}

func NewRegBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		modules: map[string]vm.Value{},
		funcs:   map[string]vm.Value{},
	}
}

func NewRegBuilderFrom(other *Registry) *RegistryBuilder {
	r := &RegistryBuilder{
		funcs:   make(map[string]vm.Value, len(other.funcs)),
		modules: make(map[string]vm.Value, len(other.modules)),
	}

	for k, f := range other.funcs {
		r.RegisterFunc(k, f.GetNativeFunction())
	}

	for k, v := range other.modules {
		r.RegisterModule(k, v.GetModule())
	}

	return r
}

func (r *RegistryBuilder) RegisterModule(name string, m map[string]vm.Value) *RegistryBuilder {
	val := vm.Value{}
	val.SetModule(m)
	r.modules[name] = val
	return r
}

func (r *RegistryBuilder) ResolveModule(name string) (vm.Value, bool) {
	val, ok := r.modules[name]
	return val, ok
}

func (r *RegistryBuilder) RemoveModule(name string) vm.Value {
	v := r.modules[name]
	delete(r.modules, name)
	return v
}

func (r *RegistryBuilder) RegisterFunc(name string, f vm.NativeFunction) *RegistryBuilder {
	val := vm.Value{}
	val.SetNativeFunction(f)
	r.funcs[name] = val
	return r
}

func (r *RegistryBuilder) ResolveFunc(name string) (vm.Value, bool) {
	val, ok := r.funcs[name]
	return val, ok
}

func (r *RegistryBuilder) RemoveFunc(name string) vm.Value {
	v := r.funcs[name]
	delete(r.funcs, name)
	return v
}

func (r *RegistryBuilder) Build() *Registry {
	funcs := make(map[string]vm.Value, len(r.funcs))
	for k, v := range r.funcs {
		funcs[k] = v
	}

	modules := make(map[string]vm.Value, len(r.modules))
	for k, v := range r.modules {
		modules[k] = v
	}

	return &Registry{
		funcs:   funcs,
		modules: modules,
	}
}

type Registry struct {
	modules map[string]vm.Value
	funcs   map[string]vm.Value
}

func (r *Registry) ResolveFunc(name string) (vm.Value, bool) {
	val, ok := r.funcs[name]
	return val, ok
}

func (r *Registry) ResolveModule(name string) (vm.Value, bool) {
	val, ok := r.modules[name]
	return val, ok
}
