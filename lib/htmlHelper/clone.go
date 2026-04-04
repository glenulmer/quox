package htmlHelper

func (e Elem_t) Clone() Elem_t {
	out := Elem_t{
		tag: e.tag,
		canWrap: e.canWrap,
		kvpairs: e.kvpairs.Copy(),
		classes: e.classes.Copy(),
		styles: e.styles.Copy(),
	}

	if len(e.wrapped) == 0 {
		return out
	}

	out.wrapped = make([]isWrappable, 0, len(e.wrapped))
	for _, item := range e.wrapped {
		switch v := item.(type) {
		case Elem_t:
			out.wrapped = append(out.wrapped, v.Clone())
		case tContent:
			out.wrapped = append(out.wrapped, tContent(string(v)))
		}
	}

	return out
}
