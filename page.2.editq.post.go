package main

import (
	"net/http"
	"strings"

	. "klpm/pkg.Global"
)

func RewriteEditQPage(w http.ResponseWriter, req *http.Request, state State_t) {
	QuoteEnsureDefaults(&state)
	vars := state.quote
	layout := RequestLayout(req)
	body := EditQBodyView(layout, vars, false)
	SendResponse(w, RewriteHTML(OuterHTML, `EditQFormBody`, body))
}

func Page2EditQEntry(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	QuoteEnsureDefaults(&state)
	EditQDropPreinsertedDependant(&state)
	SetState(req, state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}

func Page2EditQChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	value := req.FormValue(`value`)
	state := GetState(req)
	if EditQApply(&state, name, value) {
		SetState(req, state)
		RewriteEditQPage(w, req, state)
		return
	}

	SetState(req, state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}
