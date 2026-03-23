package main

import (
	"net/http"

	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func Page0Home(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	_ = state

	head := Head().
		CSS(Str(`/static/css/phone.quote.css?v=`, App.staticVersion)).
		Title(`Quo2`).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		Elem(`main`).Class(`page`, `customer-home`).Wrap(
			CustomerCard(), NL,
		), NL,
		head.Right(), NL,
	)
}
