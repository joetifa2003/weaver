package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/vm"
)

type Element interface {
	Render(parent Element) string
	AddClass(class string)
	WithClass(class string)
	SetClass(class string)
	RemoveClass(class string)
	SetAttr(key, value string)
	RemoveAttr(key string)
}

type BaseElement struct{}

func (b BaseElement) Render(parent Element) string { return "" }

func (b BaseElement) AddClass(class string) {}

func (b BaseElement) WithClass(class string) {}

func (b BaseElement) SetClass(class string) {}

func (b BaseElement) RemoveClass(class string) {}

func (b BaseElement) SetAttr(key, value string) {}

func (b BaseElement) RemoveAttr(key string) {}

type Tag struct {
	BaseElement
	Name     string
	Attrs    map[string]string
	Children []Element
}

func (t *Tag) Render(parent Element) string {
	t.Attrs = nil

	childBuf := strings.Builder{}
	for _, child := range t.Children {
		childBuf.WriteString(child.Render(t))
	}

	var mainBuf strings.Builder
	mainBuf.WriteString("<")
	mainBuf.WriteString(t.Name)
	for k, v := range t.Attrs {
		mainBuf.WriteString(" ")
		mainBuf.WriteString(k)
		mainBuf.WriteString("=\"")
		mainBuf.WriteString(v)
		mainBuf.WriteString("\"")
	}
	mainBuf.WriteString(">")

	mainBuf.WriteString(childBuf.String())

	mainBuf.WriteString("</")
	mainBuf.WriteString(t.Name)
	mainBuf.WriteString(">")
	return mainBuf.String()
}

func (t *Tag) WithClass(class string) {
	if t.Attrs == nil {
		t.Attrs = map[string]string{}
	}

	if val, ok := t.Attrs["class"]; ok {
		t.Attrs["class"] = val + " " + class
		return
	}

	t.Attrs["class"] = class
}

func (t *Tag) SetAttr(key, value string) {
	if t.Attrs == nil {
		t.Attrs = map[string]string{}
	}

	t.Attrs[key] = value
}

type Text struct {
	BaseElement
	Text string
}

func (t *Text) Render(parent Element) string {
	return t.Text
}

type Group struct {
	BaseElement
	Elements []Element
}

func (g *Group) Render(parent Element) string {
	buf := strings.Builder{}
	for _, el := range g.Elements {
		buf.WriteString(el.Render(parent))
	}
	return buf.String()
}

type WithClass struct {
	BaseElement
	Class string
}

func (w *WithClass) Render(parent Element) string {
	parent.WithClass(w.Class)
	return ""
}

type SetAttr struct {
	BaseElement
	Key   string
	Value string
}

func (s *SetAttr) Render(parent Element) string {
	parent.SetAttr(s.Key, s.Value)
	return ""
}

var htmlTags = [...]string{
	"html",
	"div",
	"span",
	"p",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"a",
	"img",
	"ul",
	"ol",
	"li",
	"form",
	"input",
	"button",
	"textarea",
	"form",
	"html",
	"head",
	"meta",
	"script",
	"link",
	"body",
	"header",
}

func registerHtmlModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("html", func() vm.Value {
		m := map[string]vm.Value{
			"withClass": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
				classArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return classArg
				}

				classStr := classArg.GetString()

				return vm.NewNativeObject(&WithClass{Class: classStr}, nil)
			}),
			"setAttr": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
				keyArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return keyArg
				}

				valArg, ok := args.Get(1, vm.ValueTypeString)
				if !ok {
					return valArg
				}

				return vm.NewNativeObject(&SetAttr{Key: keyArg.GetString(), Value: valArg.GetString()}, nil)
			}),
			"render": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
				elArg, ok := args.Get(0)
				if !ok {
					return elArg
				}

				el, ok := handleElement(v, elArg)
				if !ok {
					return vm.NewError("invalid element", vm.Value{})
				}

				return vm.NewString(el.Render(nil))
			}),
		}

		for _, tag := range htmlTags {
			m[tag] = vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
				children := make([]Element, 0, len(args))

				tag := &Tag{
					Name:  tag,
					Attrs: map[string]string{},
				}

				for _, arg := range args {
					el, ok := handleElement(v, arg)
					if !ok {
						return vm.NewError("invalid argument type", vm.Value{})
					}
					children = append(children, el)
				}

				tag.Children = children

				return vm.NewNativeObject(tag, nil)
			})
		}

		return vm.NewObject(m)
	})
}

func handleElement(v *vm.VM, val vm.Value) (Element, bool) {
	switch val.VType {
	case vm.ValueTypeString:
		return &Text{Text: val.GetString()}, true

	case vm.ValueTypeNativeObject:
		el, ok := vm.GetNativeObject[Element](val)
		if !ok {
			return nil, false
		}
		return el, true

	case vm.ValueTypeArray:
		arr := *val.GetArray()

		children := make([]Element, 0, len(arr))
		for _, val := range arr {
			el, ok := handleElement(v, val)
			if !ok {
				return nil, false
			}
			children = append(children, el)
		}

		return &Group{Elements: children}, true

	default:
		return nil, false
	}
}
