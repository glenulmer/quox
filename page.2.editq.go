package main

import (
	"net/http"

	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func EditQCSSPath() string {
	if App.layout == layoutDesktop {
		return Str(`/static/css/page.2.editq.desktop.css?v=`, App.staticVersion)
	}
	return Str(`/static/css/page.2.editq.phone.css?v=`, App.staticVersion)
}

func EditQBodyView(vars QuoteVars_t, sortForGet bool) Elem_t {
	if App.layout == layoutDesktop { return EditQDesktopBodyView(vars, sortForGet) }
	return EditQPhoneBodyView(vars, sortForGet)
}

func EditQFormView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`form`).
		Id(`EditQForm`).
		Class(`editq-form`).
		KV(`method`, `post`).
		KV(`action`, `/quote-edit-change`).
		Wrap(EditQBodyView(vars, sortForGet))
}

func EditQPageView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`main`).Class(`editq-page`).Wrap(
		EditQFormView(vars, sortForGet),
	)
}

func Page2EditQ(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	state.quote = QuoteVars(state)
	EditQEnsureFirstPlanSelected(&state)
	EditQEnsureDefaultDependent(&state)
	SetState(req, state)

	head := Head().
		CSS(EditQCSSPath()).
		JSTail(Str(`/static/js/page.2.editq.js?v=`, App.staticVersion)).
		Title(`Quo2`).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		EditQPageView(state.quote, true), NL,
		head.Right(), NL,
	)
}
