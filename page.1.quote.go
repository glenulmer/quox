package main

import (
	"net/http"

	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func QuoteCSSPath() string {
	if App.layout == layoutDesktop {
		return Str(`/static/css/page.1.quote.desktop.css?v=`, App.staticVersion)
	}
	return Str(`/static/css/page.1.quote.phone.css?v=`, App.staticVersion)
}

func QuotePageView(vars QuoteVars_t, plans QuotePlans_t) Elem_t {
	if App.layout == layoutDesktop { return QuoteDesktopPageView(vars, plans) }
	return QuotePhonePageView(vars, plans)
}

func Page1Quote(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	state.quote = QuoteVars(state)
	SetState(req, state)
	plans := QuotePlans(state)

	head := Head().
		CSS(QuoteCSSPath()).
		JSTail(Str(`/static/js/page.1.quote.buy.js?v=`, App.staticVersion)).
		JSTail(Str(`/static/js/page.1.quote.js?v=`, App.staticVersion)).
		Title(`Quo2`).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		QuotePageView(state.quote, plans), NL,
		head.Right(), NL,
	)
}
