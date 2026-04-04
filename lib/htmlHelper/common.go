package htmlHelper

import (
	. "quo2/lib/output"
)

func Div(list ...any) Elem_t { return Elem(`div`).Wrap(list...) }
func UL(list ...any) Elem_t { return Elem(`ul`).Wrap(list...) }
func LI(list ...any) Elem_t { return Elem(`li`).Wrap(list...) }
func Span(list ...any) Elem_t { return Elem(`span`).Wrap(list...) }
func Bold(list ...any) Elem_t { return Elem(`b`).Wrap(list...) }
func Ital(list ...any) Elem_t { return Elem(`i`).Wrap(list...) }

func Link(text, link string) Elem_t { return Elem(`a`).KV(`href`, link).Text(text) }
func Crumb(t, l string) string { return Link(t,l).Class(`navbar-item`, `has-text-link`).String() }

func P(list ...any) Elem_t { return Elem(`p`).Wrap(list...) }
func H1(list ...any) Elem_t { return Elem(`h1`).Class(`title`,`is-3`).Wrap(list...) }
func H2(list ...any) Elem_t { return Elem(`h2`).Class(`subtitle`,`is-5`).Wrap(list...) }
func H3(list ...any) Elem_t { return Elem(`h3`).Class(`subtitle`,`is-6`).Wrap(list...) }

func Table(list ...any) Elem_t { return Elem(`table`).Class(`table`).Wrap(list...) }

func THead(list ...any)Elem_t { return Elem(`thead`).Wrap(list...) }
func TBody(list ...any)Elem_t { return Elem(`tbody`).Wrap(list...) }
func TFoot(list ...any)Elem_t { return Elem(`tfoot`).Wrap(list...) }

func TR(list ...any) Elem_t { return Elem(`tr`).Wrap(list...) }

func TD(list ...any) Elem_t { return Elem(`td`).Wrap(list...) }
func TH(list ...any) Elem_t { return Elem(`th`).Wrap(list...) }

func THTR(names ...string) Elem_t {
	list := make([]Elem_t, len(names))
	for ix := range names {
		list[ix] = TH().Text(names[ix])
	}
	return TR(list)
}

func SubmitB() Elem_t { return Elem(`button`).Type(`submit`).Class(`submit_rec`) }
func UpdateB(text ...string) Elem_t { return SubmitB().Name(`update`).textOption(`Update`, text).Style(`display:none`) }
func DeleteB(text ...string) Elem_t { return SubmitB().Name(`delete`).textOption(`Delete`, text) }
func CreateB(text ...string) Elem_t { return SubmitB().Name(`create`).textOption(`Create`, text).Style(`display:none`) }

func (this Elem_t)textOption(normal string, alternate []string) Elem_t {
	if len(alternate) == 0 { return this.Text(normal) }
	return this.Text(alternate[0])
}

func ColGroup(items ...any) Elem_t {
	var wids []string
	var list []Elem_t
	for _, item := range items {
		switch v := item.(type) {

			case float32, float64, int, string:
				wids = append(wids, Str(v))

			case []float64:
				for _, wid := range v { wids = append(wids, Str(wid)) }
			case []int:
				for _, wid := range v { wids = append(wids, Str(wid)) }
			case []string:
				for _, wid := range v { wids = append(wids, Str(wid)) }
		
		}
	}
	for _, wid := range wids {
		list = append(list, Elem(`col`).Style(Str(`width:`, wid, `em`)))
	}
	return Elem(`colgroup`).Wrap(list)
}
