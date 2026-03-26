package main

import . "pm/lib/htmlHelper"

func QuoteDesktopFormBodyView(vars QuoteVars_t) Elem_t {
	return Div().Id(`QuoteFormBody`).Wrap(
		Div().Class(`quote-desk-workbench`).Wrap(
			Div().Class(`quote-desk-rows`).Wrap(
				Div().Class(`quote-desk-row`, `quote-desk-row-top`).Wrap(
					QuoteNamedControlOnlySpanView(layoutDesktop, `custName`, vars, 0, `quote-desk-no-label`, `quote-desk-name`),
					QuoteNamedControlOnlySpanView(layoutDesktop, `segment`, vars, 0, `quote-desk-no-label`, `quote-desk-segment`),
					Div().Class(`quote-desk-flags`).Wrap(
						QuoteNamedControlSpanView(layoutDesktop, `vision`, vars, 1, `quote-desk-flag`),
						QuoteNamedControlSpanView(layoutDesktop, `tempVisa`, vars, 1, `quote-desk-flag`),
						QuoteNamedControlSpanView(layoutDesktop, `noPVN`, vars, 1, `quote-desk-flag`),
						QuoteNamedControlSpanView(layoutDesktop, `naturalMed`, vars, 1, `quote-desk-flag`),
					),
				),
				Div().Class(`quote-desk-row`, `quote-desk-row-mid`).Wrap(
					QuoteNamedControlSpanView(layoutDesktop, `birth`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlSpanView(layoutDesktop, `buy`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlSpanView(layoutDesktop, `sickCover`, vars, 0, `quote-desk-compact`, `quote-desk-right`),
					QuoteNamedControlSpanView(layoutDesktop, `priorCov`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlSpanView(layoutDesktop, `exam`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlSpanView(layoutDesktop, `specref`, vars, 0, `quote-desk-compact`),
				),
				Div().Class(`quote-desk-row`, `quote-desk-row-bottom`).Wrap(
					QuoteNamedControlLabelSpanView(layoutDesktop, `deductibleMin`, `Deductible Min`, vars, 0, `quote-desk-compact`, `quote-desk-right`),
					QuoteNamedControlLabelSpanView(layoutDesktop, `deductibleMax`, `Max`, vars, 0, `quote-desk-compact`, `quote-desk-right`),
					QuoteNamedControlLabelSpanView(layoutDesktop, `hospitalMin`, `Hospital Min`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlLabelSpanView(layoutDesktop, `hospitalMax`, `Max`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlLabelSpanView(layoutDesktop, `dentalMin`, `Dental Min`, vars, 0, `quote-desk-compact`),
					QuoteNamedControlLabelSpanView(layoutDesktop, `dentalMax`, `Max`, vars, 0, `quote-desk-compact`),
				),
			),
		),
		QuoteDesktopSelectedPlansBox(vars),
	)
}

func QuoteDesktopFormView(vars QuoteVars_t) Elem_t {
	return Elem(`form`).
		Id(`QuoteForm`).
		Class(`quote-form`, `quote-form-desktop`).
		KV(`method`, `post`).
		KV(`action`, `/quote-info-change`).
		Wrap(QuoteDesktopFormBodyView(vars))
}

func QuoteDesktopPlansView(data QuotePlans_t) Elem_t {
	return Div().Id(`QuotePlans`).Class(`quote-plan-results`, `quote-desktop-results`).Wrap(
		QuotePlanDesktopView(data),
		QuoteFilteredPlansBox(data.filtered),
	)
}

func QuoteDesktopPageView(vars QuoteVars_t, plans QuotePlans_t) Elem_t {
	return Elem(`main`).Class(`quote-page`, `quote-page-desktop`).Wrap(
		QuoteDesktopFormView(vars),
		QuoteDesktopPlansView(plans),
	)
}
