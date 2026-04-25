package main

import . "klpm/lib/htmlHelper"

func EditQPhoneBodyView(vars UIBagVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`, `editq-body-phone`).Wrap(
		EditQHeaderView(vars),
		EditQDependentsView(vars, sortForGet),
		EditQQuoteReviewCardView(vars),
	)
}
