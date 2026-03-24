package main

import (
	"net/http"

	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func Page0Home(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)

	head := Head().
		CSS(Str(`/static/css/phone.quote.css?v=`, App.staticVersion)).
		Title(`Quo2`).
		End()

	card := CustomerCard()
	state.quote = card.AllValues()
	SetState(req, state)

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		Elem(`main`).Class(`page`, `customer-home`).Wrap(
			card, NL,
		), NL,
		head.Right(), NL,
	)
}
