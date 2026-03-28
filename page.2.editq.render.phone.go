package main

import . "pm/lib/htmlHelper"

func EditQPhoneBodyView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`, `editq-body-phone`).Wrap(
		EditQHeaderView(vars),
		EditQPrimeChargesView(vars),
		EditQDependentsView(vars, sortForGet),
	)
}
