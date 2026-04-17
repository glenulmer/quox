package main

import . "quo2/lib/htmlHelper"

func QuotePhoneFormBodyView(vars UIBagVars_t) Elem_t {
	return Div().Id(`QuoteFormBody`).Wrap(
		Div().Id(`QuotePhoneStickyAnchor`).Class(`quote-phone-sticky-anchor`),
		Elem(`details`).Id(`QuoteInfoCard`).Class(`quote-card`, `quote-phone-card`, `quote-phone-fold`, `quote-phone-info-fold`).KV(`open`, `open`).Wrap(
			Elem(`summary`).Class(`quote-card-title`, `quote-phone-fold-title`, `quote-phone-selected-title`).Wrap(
				Span(`Quote Info`),
				QuoteResetButton(`quote-phone-selected-edit-btn`),
			),
			Div().Class(`quote-grid`, `quote-grid-phone`).Wrap(
			QuoteNamedControlOnlySpanView(layoutPhone, `clientName`, vars, 8, `quote-phone-no-label`),
			QuoteNamedControlOnlySpanView(layoutPhone, `segment`, vars, 4, `quote-phone-no-label`),
			QuoteNamedControlSpanView(layoutPhone, `birth`, vars, 4),
			QuoteNamedControlSpanView(layoutPhone, `buy`, vars, 4),
			QuoteNamedControlSpanView(layoutPhone, `sickCover`, vars, 4, `quote-phone-right`),
			QuoteNamedControlSpanView(layoutPhone, `priorCov`, vars, 4),
			QuoteNamedControlSpanView(layoutPhone, `exam`, vars, 4),
			QuoteNamedControlSpanView(layoutPhone, `specref`, vars, 4),

			QuoteSpacer(1),
			Div().Class(QuoteSpanClass(10), `quote-phone-checks`).Wrap(
				QuoteNamedControlSpanView(layoutPhone, `vision`, vars, 1, `quote-phone-check`),
				QuoteNamedControlSpanView(layoutPhone, `tempVisa`, vars, 1, `quote-phone-check`),
				QuoteNamedControlSpanView(layoutPhone, `noPVN`, vars, 1, `quote-phone-check`),
				QuoteNamedControlSpanView(layoutPhone, `naturalMed`, vars, 1, `quote-phone-check`),
			),
			QuoteSpacer(1),

			QuoteNamedControlLabelSpanView(layoutPhone, `deductibleMin`, `Deductible`, vars, 4, `quote-phone-right`),
			QuoteNamedControlLabelSpanView(layoutPhone, `hospitalMin`, `Hospital`, vars, 4),
			QuoteNamedControlLabelSpanView(layoutPhone, `dentalMin`, `Dental`, vars, 4),

			QuoteNamedControlOnlySpanView(layoutPhone, `deductibleMax`, vars, 4, `quote-phone-no-label`, `quote-phone-right`),
			QuoteNamedControlOnlySpanView(layoutPhone, `hospitalMax`, vars, 4, `quote-phone-no-label`),
			QuoteNamedControlOnlySpanView(layoutPhone, `dentalMax`, vars, 4, `quote-phone-no-label`),
			),
		),
		QuotePhoneSelectedPlansBox(vars),
	)
}

func QuotePhoneFormView(vars UIBagVars_t) Elem_t {
	return Elem(`form`).
		Id(`QuoteForm`).
		Class(`quote-form`, `quote-form-phone`).
		KV(`method`, `post`).
		KV(`action`, `/quote-info-change`).
		Wrap(QuotePhoneFormBodyView(vars))
}

func QuotePhonePlansView(data QuotePlans_t) Elem_t {
	var plans []Elem_t
	for _, x := range data.plans { plans = append(plans, QuotePlanCardView(x)) }
	return Div().Id(`QuotePlans`).Class(`quote-plan-results`, `quote-phone-results`).Wrap(
		Elem(`details`).Class(`quote-phone-fold`, `quote-phone-plans-fold`).KV(`open`, `open`).Wrap(
			Elem(`summary`).Class(`quote-plan-toolbar`, `quote-phone-plan-toolbar`, `quote-phone-fold-title`).Wrap(
				Div(`Plans (` , len(data.plans), `)`).Class(`quote-card-title`),
				Div().Class(`quote-plan-sort`).Wrap(
					QuoteSortSelectView(data.sortBy),
				),
			),
			Div().Class(`quote-phone-plans-body`).Wrap(
				Div().Class(`quote-plan-list`, `quote-plan-list-phone`).Wrap(plans),
			),
		),
		QuoteFilteredPlansBox(data.filtered),
	)
}

func QuotePhonePageView(vars UIBagVars_t, plans QuotePlans_t) Elem_t {
	return Elem(`main`).Class(`quote-page`, `quote-page-phone`).Wrap(
		QuotePhoneFormView(vars),
		QuotePhonePlansView(plans),
	)
}
