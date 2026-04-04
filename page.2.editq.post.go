package main

import (
	"net/http"
	"strings"

	. "pm/pkg.Global"
)

func RewriteEditQPage(w http.ResponseWriter, state State_t) {
	vars := QuoteVars(state)
	body := EditQBodyView(vars, false)
	SendResponse(w, RewriteHTML(OuterHTML, `EditQFormBody`, body))
}

func Page2EditQEntry(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	state.quote = QuoteVars(state)
	EditQEnsureDefaultDependent(&state)
	SetState(req, state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}

func Page2EditQChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	value := req.FormValue(`value`)
	state := GetState(req)
	if EditQApply(&state, name, value) {
		SetState(req, state)
		RewriteEditQPage(w, state)
		return
	}

	SetState(req, state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}
