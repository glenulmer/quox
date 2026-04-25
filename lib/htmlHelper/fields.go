package htmlHelper

import (
	. `klpm/lib/output`
)

func wid12(wid int) string {
	switch {
		case wid < 1: wid = 12;
		case wid > 12: wid = 12;
	}
	return Str(`w12-`, wid)
}

func Label(x any) Elem_t { return Elem(`label`).Wrap(x) }

func Field(items ...any) Elem_t {
	d := Div().Class(`card-field`)
	var wid, label string
	var hasInt, hasLabel bool

	for _, item := range items {
		if hasInt && hasLabel { break }
		switch it := item.(type) {
		case int:
			if hasInt { continue }
			wid = wid12(it)
			hasInt = true
		case string:
			if hasLabel { continue }
			label = it
			hasLabel = true
		}
	}

	if label != `` { d = d.Wrap(Label(label)) }
	if wid != `` { d = d.Class(wid) }

	return d
}

func DateInput() Elem_t { return Elem(`input`).KV(`autocomplete`,`off`).KV(`type`,`date`) }

func (in Elem_t)Req() Elem_t { return in.KV(`required`) }

func (in Elem_t)AllValues() map[string]string {
	out := make(map[string]string)
	in.addAllValues(out)
	return out
}

func (in Elem_t)addAllValues(out map[string]string) {
	name := in.attr(`name`)
	if name != `` {
		switch in.tag {
		case `input`:
			out[name] = in.attr(`value`)
		case `select`:
			out[name] = in.selectValue()
		}
	}

	for _, wrapped := range in.wrapped {
		if elem, ok := wrapped.(Elem_t); ok { elem.addAllValues(out) }
	}
}

func (in Elem_t)attr(name string) string {
	val, _ := in.kvpairs.data[name]
	return val
}

func (in Elem_t)selectValue() string {
	type tOptionValue struct {
		value string
		selected bool
	}

	var first string
	var hasFirst bool

	var scan func(Elem_t) (tOptionValue, bool)
	scan = func(e Elem_t) (tOptionValue, bool) {
		if e.tag == `option` {
			val := e.attr(`value`)
			_, selected := e.kvpairs.data[`selected`]
			return tOptionValue{ value: val, selected: selected }, true
		}

		for _, wrapped := range e.wrapped {
			child, ok := wrapped.(Elem_t)
			if !ok { continue }
			option, found := scan(child)
			if !found { continue }
			if !hasFirst {
				first = option.value
				hasFirst = true
			}
			if option.selected {
				return option, true
			}
		}
		return tOptionValue{}, false
	}

	if option, found := scan(in); found && option.selected {
		return option.value
	}
	return first
}
