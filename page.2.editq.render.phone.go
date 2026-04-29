package main

import . "klpm/lib/htmlHelper"

func EditQPhoneBodyView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`, `editq-body-phone`).Wrap(
		EditQHeaderView(vars),
		EditQDependantsView(vars, sortForGet),
		EditQQuoteReviewCardView(vars),
	)
}
