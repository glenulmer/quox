package main

import . "pm/lib/htmlHelper"

func EditQDesktopBodyView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`, `editq-body-desktop`).Wrap(
		EditQHeaderView(vars),
		Div().Class(`editq-main`).Wrap(
			EditQPrimeChargesView(vars),
			EditQDependentsView(vars, sortForGet),
		),
	)
}
