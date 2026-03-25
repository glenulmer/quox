package main

import (
	"net/http"
	"strings"

	. "pm/pkg.Global"
	. "pm/lib/output"
)

func RewritePlansAndFilters(w http.ResponseWriter, state State_t) {
	SendResponse(w,
		RewriteHTML(OuterHTML, `PlansAndFilters`, ListPlans(state)),
	)
}

func QuoteInfoChangeHandler(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	value := req.FormValue(`value`)

	state := GetState(req)
	if state.quote == nil {
		state.quote = make(QuoteVars_t)
	}
	if name != `` {
		state.quote[name] = Q(value)
		SetState(req, state)
	}

	RewritePlansAndFilters(w, state)
}
