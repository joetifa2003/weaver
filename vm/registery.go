package vm

import (
	"github.com/joetifa2003/weaver/internal/pkg/ds"
)

type RegistryBuilder struct {
	modules *ds.ConcMap[string, Value]
	funcs   *ds.ConcMap[string, Value]
}

func NewRegBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		modules: ds.NewConcMap[string, Value](),
		funcs:   ds.NewConcMap[string, Value](),
	}
}

func NewRegBuilderFrom(other *Registry) *RegistryBuilder {
	r := &RegistryBuilder{
		funcs:   ds.NewConcMap[string, Value](),
		modules: ds.NewConcMap[string, Value](),
	}

	for k, f := range other.funcs.Iter() {
		r.RegisterFunc(k, f.GetNativeFunction())
	}

	for k, v := range other.modules.Iter() {
		r.RegisterModule(k, v)
	}

	return r
}

func (r *RegistryBuilder) RegisterModule(name string, v Value) *RegistryBuilder {
	r.modules.Set(name, v)
	return r
}

func (r *RegistryBuilder) ResolveModule(name string) (Value, bool) {
	return r.modules.Get(name)
}

func (r *RegistryBuilder) RemoveModule(name string) Value {
	v, ok := r.modules.Get(name)
	if !ok {
		return Value{}
	}
	r.modules.Delete(name)
	return v
}

func (r *RegistryBuilder) RegisterFunc(name string, f NativeFunction) *RegistryBuilder {
	val := Value{}
	val.SetNativeFunction(f)
	r.funcs.Set(name, val)
	return r
}

func (r *RegistryBuilder) ResolveFunc(name string) (Value, bool) {
	return r.funcs.Get(name)
}

func (r *RegistryBuilder) RemoveFunc(name string) Value {
	v, ok := r.funcs.Get(name)
	if !ok {
		return Value{}
	}
	r.funcs.Delete(name)
	return v
}

func (r *RegistryBuilder) Build() *Registry {
	funcs := ds.NewConcMap[string, Value]()
	for k, v := range r.funcs.Iter() {
		funcs.Set(k, v)
	}

	modules := ds.NewConcMap[string, Value]()
	for k, v := range r.modules.Iter() {
		modules.Set(k, v)
	}

	return &Registry{
		funcs:   funcs,
		modules: modules,
	}
}

type Registry struct {
	modules *ds.ConcMap[string, Value]
	funcs   *ds.ConcMap[string, Value]
}

func (r *Registry) ResolveFunc(name string) (Value, bool) {
	return r.funcs.Get(name)
}

func (r *Registry) ResolveModule(name string) (Value, bool) {
	return r.modules.Get(name)
}
