package main

import (
	"net/http"
	"strings"

	. "pm/pkg.Global"
)

func RewriteQuotePage(w http.ResponseWriter, state State_t) {
	vars := QuoteVars(state)
	plans := QuotePlans(state)
	form := any(QuotePhoneFormView(vars))
	planView := any(QuotePhonePlansView(plans))
	if App.layout == layoutDesktop {
		form = QuoteDesktopFormView(vars)
		planView = QuoteDesktopPlansView(plans)
	}
	SendResponse(w,
		RewriteHTML(OuterHTML, `QuoteForm`, form),
		RewriteHTML(OuterHTML, `QuotePlans`, planView),
	)
}

func Page1QuoteChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	state := GetState(req)
	if name != `` {
		QuoteApply(&state, name, req.FormValue(`value`))
		SetState(req, state)
		RewriteQuotePage(w, state)
		return
	}

	QuoteApplyForm(&state, req)
	SetState(req, state)
	http.Redirect(w, req, `/`, http.StatusSeeOther)
}
