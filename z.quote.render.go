package main

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func QuoteVarBool(v string) bool {
	switch Lower(Trim(v)) {
	case `1`, `on`, `yes`, `true`:
		return true
	}
	return false
}

func QuoteSelect(name, value string, choices []QuoteChoice_t) Elem_t {
	var options []Elem_t
	for _, x := range choices {
		opt := Option().KV(`value`, x.id).Text(x.label)
		if value == Str(x.id) { opt = opt.KV(`selected`) }
		options = append(options, opt)
	}
	return Select(options).Name(name)
}

func QuoteCheckbox(name string, checked bool) Elem_t {
	in := Elem(`input`).Type(`checkbox`).Name(name).KV(`value`, `1`)
	if checked { in = in.KV(`checked`) }
	return in
}

func QuoteInputText(name, value, placeholder string) Elem_t {
	in := Elem(`input`).Type(`text`).Name(name).Value(value).KV(`autocomplete`, `off`)
	if placeholder != `` { in = in.Place(placeholder) }
	return in
}

func QuoteInputDate(name, value string) Elem_t {
	return DateInput().
		Name(name).
		Value(value).
		KV(`inputmode`, `none`).
		KV(`onclick`, `if(this.showPicker){this.showPicker();}`).
		KV(`onfocus`, `if(this.showPicker){this.showPicker();}`).
		KV(`onkeydown`, `return false`).
		KV(`onpaste`, `return false`).
		KV(`ondrop`, `return false`)
}

func QuoteInputNumber(name, value string, min, max, step int) Elem_t {
	in := Elem(`input`).Type(`number`).Name(name).Value(value)
	if min != 0 || max != 0 || step != 0 {
		in = in.KV(`min`, min).KV(`max`, max).KV(`step`, step)
	}
	return in
}

func QuoteControlInput(x QuoteControl_t, vars QuoteVars_t) Elem_t {
	value := vars[x.name]
	switch x.kind {
	case quoteText:
		return QuoteInputText(x.name, value, x.placeholder)
	case quoteDate:
		return QuoteInputDate(x.name, value)
	case quoteNumber:
		return QuoteInputNumber(x.name, value, x.min, x.max, x.step)
	case quoteSelect:
		return QuoteSelect(x.name, value, QuoteControlChoices(x))
	case quoteCheckbox:
		return QuoteCheckbox(x.name, QuoteVarBool(value))
	}
	return Span(`missing control`)
}

func QuoteSpanClass(span int) string {
	if span < 1 { span = 1 }
	if span > 12 { span = 12 }
	return Str(`quote-span-`, span)
}

func QuoteControlViewSpan(layout string, x QuoteControl_t, vars QuoteVars_t, span int, class ...string) Elem_t {
	if span == 0 { span = QuoteControlSpan(x, layout) }
	classes := append([]string{ QuoteSpanClass(span) }, class...)
	if x.kind == quoteCheckbox {
		return Elem(`label`).Class(`quote-check`).Class(classes...).Wrap(
			QuoteControlInput(x, vars),
			Span(x.label).Class(`quote-check-text`),
		)
	}
	return Elem(`label`).Class(`quote-field`).Class(classes...).Wrap(
		Span(x.label).Class(`quote-label`),
		QuoteControlInput(x, vars),
	)
}

func QuoteControlView(layout string, x QuoteControl_t, vars QuoteVars_t) Elem_t {
	return QuoteControlViewSpan(layout, x, vars, 0)
}

func QuoteControlLabelViewSpan(layout string, x QuoteControl_t, label string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x.label = label
	return QuoteControlViewSpan(layout, x, vars, span, class...)
}

func QuoteNamedControlLabelSpanView(layout, name, label string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlLabelViewSpan(layout, x, label, vars, span, class...)
}

func QuoteControlOnlyViewSpan(layout string, x QuoteControl_t, vars QuoteVars_t, span int, class ...string) Elem_t {
	if span == 0 { span = QuoteControlSpan(x, layout) }
	classes := append([]string{ QuoteSpanClass(span) }, class...)
	return Div().Class(`quote-control`).Class(classes...).Wrap(
		QuoteControlInput(x, vars),
	)
}

func QuoteNamedControlOnlySpanView(layout, name string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlOnlyViewSpan(layout, x, vars, span, class...)
}

func QuoteSpacer(span int, class ...string) Elem_t {
	classes := append([]string{ QuoteSpanClass(span), `quote-spacer` }, class...)
	return Div().Class(classes...)
}

func QuoteControlViews(layout, group string, vars QuoteVars_t) []Elem_t {
	var out []Elem_t
	for _, x := range QuoteControlsByGroup(layout, group) {
		out = append(out, QuoteControlView(layout, x, vars))
	}
	return out
}

func QuoteNamedControlView(layout, name string, vars QuoteVars_t) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlView(layout, x, vars)
}

func QuoteNamedControlSpanView(layout, name string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlViewSpan(layout, x, vars, span, class...)
}
