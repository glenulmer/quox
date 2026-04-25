package main

import . "klpm/lib/htmlHelper"

func EditQDesktopBodyView(vars UIBagVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`, `editq-body-desktop`).Wrap(
		EditQHeaderView(vars),
		EditQDependentsView(vars, sortForGet),
		EditQQuoteReviewCardView(vars),
	)
}
