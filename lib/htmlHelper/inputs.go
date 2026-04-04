package htmlHelper

import . `quo2/lib/output`

func Checkbox(name string, value bool) Elem_t {
	return Elem(`input`).Type(`checkbox`).Name(name).CO(value)
}

func Check(ch bool) string { if !ch { return `` }; return ` checked` }

func (e Elem_t)Check(ch bool) Elem_t { if !ch { return e }; return e.KV(`checked`) }

func Wedge() Elem_t { return Div().Class(`wedge`) }
func WedgeC() Elem_t { return Div().Class(`wedgeC`) }
func TextIn() Elem_t { return Elem(`input`).KV(`autocomplete`,`off`).KV(`type`,`text`).Class(`is-small`, `kledit`) }

func CrudText(name string, value any) Elem_t {
	return TextIn().Name(name).VO(value)
}

func CrudCheck(name string, value bool) Elem_t {
	return CBox(name, value)
}

func CrudSelect(name string, selected any, options ...Elem_t) Elem_t {
	return Select(options).Name(name).SelO(selected)
}

func Select(items ...any) Elem_t { return Elem(`select`).Wrap(items...) }

func (e Elem_t)Choose(newVal any) Elem_t {
	if e.tag != `select` { return e }
	xVal := Q(Str(newVal))
	chosen := false
	for ix, wrapped := range e.wrapped {
		if elem, ok := wrapped.(Elem_t); ok {
			val, _ := elem.kvpairs.data[`value`]
			if !chosen && val == xVal {
				e.wrapped[ix] = elem.KV(`selected`)
				chosen = true
				continue
			}
			e.wrapped[ix] = elem.CutAttrib(`selected`)
		}
	}
	return e
}

func (e Elem_t)GetSelected() string {
	if e.tag != `select` { return `` }
	for _, wrapped := range e.wrapped {
		if elem, ok := wrapped.(Elem_t); ok {
			_, ok := elem.kvpairs.data[`selected`]
			if ok {
				val, _ := elem.kvpairs.data[`value`]
				return val
			}
		}
	}
	return ``
}

func Option() Elem_t { return Elem(`option`) }
/*
<select>
	<option value="1">Yes</option>
	<option value="0">No</option>
</select>
*/
