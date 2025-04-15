package registry

import (
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/vm"
)

type RegistryBuilder struct {
	modules *ds.ConcMap[string, vm.Value]
	funcs   *ds.ConcMap[string, vm.Value]
}

func NewRegBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		modules: ds.NewConcMap[string, vm.Value](),
		funcs:   ds.NewConcMap[string, vm.Value](),
	}
}

func NewRegBuilderFrom(other *Registry) *RegistryBuilder {
	r := &RegistryBuilder{
		funcs:   ds.NewConcMap[string, vm.Value](),
		modules: ds.NewConcMap[string, vm.Value](),
	}

	for k, f := range other.funcs.Iter() {
		r.RegisterFunc(k, f.GetNativeFunction())
	}

	for k, v := range other.modules.Iter() {
		r.RegisterModule(k, v.GetModule())
	}

	return r
}

func (r *RegistryBuilder) RegisterModule(name string, m map[string]vm.Value) *RegistryBuilder {
	val := vm.Value{}
	val.SetModule(m)
	r.modules.Set(name, val)
	return r
}

func (r *RegistryBuilder) ResolveModule(name string) (vm.Value, bool) {
	return r.modules.Get(name)
}

func (r *RegistryBuilder) RemoveModule(name string) vm.Value {
	v, ok := r.modules.Get(name)
	if !ok {
		return vm.Value{}
	}
	r.modules.Delete(name)
	return v
}

func (r *RegistryBuilder) RegisterFunc(name string, f vm.NativeFunction) *RegistryBuilder {
	val := vm.Value{}
	val.SetNativeFunction(f)
	r.funcs.Set(name, val)
	return r
}

func (r *RegistryBuilder) ResolveFunc(name string) (vm.Value, bool) {
	return r.funcs.Get(name)
}

func (r *RegistryBuilder) RemoveFunc(name string) vm.Value {
	v, ok := r.funcs.Get(name)
	if !ok {
		return vm.Value{}
	}
	r.funcs.Delete(name)
	return v
}

func (r *RegistryBuilder) Build() *Registry {
	funcs := ds.NewConcMap[string, vm.Value]()
	for k, v := range r.funcs.Iter() {
		funcs.Set(k, v)
	}

	modules := ds.NewConcMap[string, vm.Value]()
	for k, v := range r.modules.Iter() {
		modules.Set(k, v)
	}

	return &Registry{
		funcs:   funcs,
		modules: modules,
	}
}

type Registry struct {
	modules *ds.ConcMap[string, vm.Value]
	funcs   *ds.ConcMap[string, vm.Value]
}

func (r *Registry) ResolveFunc(name string) (vm.Value, bool) {
	return r.funcs.Get(name)
}

func (r *Registry) ResolveModule(name string) (vm.Value, bool) {
	return r.modules.Get(name)
}
