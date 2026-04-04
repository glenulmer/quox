package htmlHelper

type tCB struct {
	name string
	value bool
}

func CB(name string, value bool) tCB { return tCB{ name:name, value:value } }

func (in Elem_t)Tag(s string) Elem_t { return in.KV(`data-tag`, s) }

func inSpan(tag string) Elem_t { return Elem(`span`).Class(`inputSpan`).Tag(tag) }

func CBox(name string, value bool) Elem_t { return inSpan(name).Wrap(Checkbox(name, value)).Tag(name) }

func CBoxes(tag string, pairs ...tCB) Elem_t {
	span := inSpan(tag)
	boxes := make([]Elem_t, len(pairs))
	for i, pair := range pairs {
		boxes[i] = Checkbox(pair.name, pair.value)
	}
	return span.Wrap(boxes)
}
