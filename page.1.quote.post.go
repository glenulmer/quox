package main

import (
	"net/http"
	"strings"

	. "pm/pkg.Global"
)

func RewriteQuotePage(w http.ResponseWriter, state State_t) {
	vars := QuoteVars(state)
	plans := QuotePlans(state)
	form := any(QuotePhoneFormBodyView(vars))
	planView := any(QuotePhonePlansView(plans))
	if App.layout == layoutDesktop {
		form = QuoteDesktopFormBodyView(vars)
		planView = QuoteDesktopPlansView(plans)
	}
	SendResponse(w,
		RewriteHTML(OuterHTML, `QuoteFormBody`, form),
		RewriteHTML(OuterHTML, `QuotePlans`, planView),
	)
}

func Page1QuoteChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	state := GetState(req)
	if name == QuoteResetControlName() {
		state.quote = QuoteDefaultVars()
		SetState(req, state)
		RewriteQuotePage(w, state)
		return
	}
	if QuoteSelectedApply(&state, name, req.FormValue(`value`)) {
		SetState(req, state)
		RewriteQuotePage(w, state)
		return
	}
	QuoteApply(&state, name, req.FormValue(`value`))
	SetState(req, state)
	RewriteQuotePage(w, state)
}
