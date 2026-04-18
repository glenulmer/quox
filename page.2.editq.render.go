package main

import (
	. "quo2/lib/dec2"
	. "quo2/lib/htmlHelper"
	. "quo2/lib/output"
)

func EditQDelButton(name string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(name).
		Value(`1`).
		Class(`editq-del-btn`).
		Text(`🗑`)
}

func EditQAddButton(name, label string, class ...string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(name).
		Value(`1`).
		Class(`editq-add-btn`).
		Class(class...).
		Text(label)
}

func EditQPreexChargeModeView(name, mode string) Elem_t {
	return Select(
		Option().KV(`value`, editQPreexModePct).Text(`%`),
		Option().KV(`value`, editQPreexModeEur).Text(`€`),
	).Name(name).Choose(EditQPreexMode(mode)).Class(`editq-preex-mode`)
}

func EditQPreexChargeRowView(x EditQPreexCharge_t, modeKey, amountKey, noteKey string) Elem_t {
	return Div().Class(`editq-preex-row`).Wrap(
		Div(x.level).Class(`editq-preex-categ`),
		Div().Class(`editq-preex-inputs`).Wrap(
			EditQPreexChargeModeView(modeKey, x.mode),
			QuoteInputText(amountKey, x.amount, `Amount`).Class(`editq-preex-amount`, `editq-preex-amount-input`),
			QuoteInputText(noteKey, x.note, `Optional note`).Class(`editq-preex-note`),
		),
	)
}

func EditQEuroCentText(v EuroCent_t) string {
	return Str(v.String(), ` €`)
}

func EditQPreexTotal(vars UIBagVars_t) EuroCent_t {
	charges := EditQPreexCharges(vars)
	var total EuroCent_t
	for _, x := range charges { total += x.applied }
	return total
}

func EditQPreexTitleText(vars UIBagVars_t) string {
	totalFlat := EuroFlatFromCent(EditQPreexTotal(vars))
	return Str(`Pre-existing conditions (`, totalFlat.OutEuro(), `)`)
}

func EditQPreexChargesView(vars UIBagVars_t) Elem_t {
	charges := EditQPreexCharges(vars)
	appliedByItem := make(map[int]EuroCent_t)
	planByItem := make(map[int]EuroCent_t)
	for _, x := range charges {
		appliedByItem[x.itemId] += x.applied
		if _, ok := planByItem[x.itemId]; !ok { planByItem[x.itemId] = x.planPrice }
	}

	var rows []Elem_t
	prevItemId := 0
	hasPlan := false
	for _, x := range charges {
		if x.itemId != prevItemId {
			if hasPlan {
				base := planByItem[prevItemId]
				preex := appliedByItem[prevItemId]
				rows = append(rows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
					Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
						Div().Class(`editq-preex-summary-spacer`),
						Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
						Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
					),
				))
			}
			rows = append(rows, Div().Class(`editq-preex-plan`).Wrap(
				Span(x.plan).Class(`editq-preex-plan-name`),
			))
			prevItemId = x.itemId
			hasPlan = true
		}
		rows = append(rows, EditQPreexChargeRowView(x, EditQPreexModeKey(x.itemId, x.categId), EditQPreexAmountKey(x.itemId, x.categId), EditQPreexNoteKey(x.itemId, x.categId)))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		preex := appliedByItem[prevItemId]
		rows = append(rows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
			Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
				Div().Class(`editq-preex-summary-spacer`),
				Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
				Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
			),
		))
	}
	if len(rows) == 0 {
		rows = append(rows, Div(`No payable selected plan categories yet.`).Class(`editq-preex-empty`))
	}
	return Div().Class(`editq-section`, `editq-preex-charges`).Wrap(
		Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
		Div().Class(`editq-preex-list`).Wrap(rows),
	)
}

func EditQDependentView(vars UIBagVars_t, dep EditQDep_t) Elem_t {
	charges := EditQDependentCharges(vars, dep)
	appliedByItem := make(map[int]EuroCent_t)
	planByItem := make(map[int]EuroCent_t)
	for _, x := range charges {
		appliedByItem[x.itemId] += x.applied
		if _, ok := planByItem[x.itemId]; !ok { planByItem[x.itemId] = x.planPrice }
	}

	var depRows []Elem_t
	prevItemId := 0
	hasPlan := false
	for _, x := range charges {
		if x.itemId != prevItemId {
			if hasPlan {
				base := planByItem[prevItemId]
				preex := appliedByItem[prevItemId]
				depRows = append(depRows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
					Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
						Div().Class(`editq-preex-summary-spacer`),
						Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
						Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
					),
				))
			}
			ageText := ``
			if x.planAge > 0 {
				ageText = Str(`Effective age: `, x.planAge)
				if x.planAgeMode == `exact` { ageText = Str(`Effective age (Exact): `, x.planAge) }
				if x.planAgeMode == `year` { ageText = Str(`Effective age (Year): `, x.planAge) }
			}
			if ageText == `` {
				depRows = append(depRows, Div().Class(`editq-preex-plan`).Wrap(
					Span(x.plan).Class(`editq-preex-plan-name`),
				))
			} else {
				depRows = append(depRows, Div().Class(`editq-preex-plan`).Wrap(
					Div().Class(`editq-preex-plan-head`).Wrap(
						Span(x.plan).Class(`editq-preex-plan-name`),
					),
					Div(ageText).Class(`editq-preex-plan-age`),
				))
			}
			prevItemId = x.itemId
			hasPlan = true
		}
		depRows = append(depRows, EditQPreexChargeRowView(
			x,
			EditQDepChargeModeKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeAmountKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeNoteKey(dep.depId, x.itemId, x.categId),
		))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		preex := appliedByItem[prevItemId]
		depRows = append(depRows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
			Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
				Div().Class(`editq-preex-summary-spacer`),
				Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
				Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
			),
		))
	}
	if len(depRows) == 0 {
		depRows = append(depRows, Div(`No selected plan available.`).Class(`editq-preex-empty`))
	}
	return Div().Class(`editq-dependent`).Wrap(
		Div().Class(`editq-dependent-fields-row`).Wrap(
			QuoteInputText(EditQDepNameKey(dep.depId), dep.name, `Dependent name`).Class(`editq-dep-name`),
			QuoteInputDate(EditQDepBirthKey(dep.depId), dep.birth).Class(`editq-dep-birth`),
			Elem(`label`).Class(`editq-check`, `editq-dep-vision`).KV(`title`, `Vision`).Wrap(
				QuoteCheckbox(EditQDepVisionKey(dep.depId), dep.vision),
				Span(`Vision`).Class(`editq-check-text`),
			),
			EditQDelButton(EditQDepDelControlName(dep.depId)).Class(`editq-dep-del`),
		),
		Div().Class(`editq-preex-charges`, `editq-dependent-preex`).Wrap(
			Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
			Div().Class(`editq-preex-list`).Wrap(depRows),
		),
	)
}

func EditQDependentsView(vars UIBagVars_t, sortForGet bool) Elem_t {
	deps := EditQDependents(vars, sortForGet)
	namedCount := 0
	var list []Elem_t
	for _, dep := range deps {
		if Trim(dep.name) != `` { namedCount++ }
		list = append(list, EditQDependentView(vars, dep))
	}
	if len(deps) < editQDepMaxCount {
		list = append(list, EditQAddButton(EditQDepAddControlName(), `Add`, `editq-add-dependent`))
	}
	return EditQTopCardView(
		`EditQDependentsCard`,
		Str(`Dependents (`, namedCount, `)`),
		false,
		Div().Class(`editq-section`, `editq-dependents`).Wrap(
			Div().Class(`editq-dependent-list`).Wrap(list),
		),
	)
}

func EditQTopCardView(id, title string, open bool, body Elem_t, right ...Elem_t) Elem_t {
	var titleRow []Elem_t
	titleRow = append(titleRow, Span(title).Class(`editq-card-title-text`))
	if len(right) > 0 {
		titleRow = append(titleRow, Div().Class(`editq-card-title-right`).Wrap(right))
	}
	card := Elem(`details`).Id(id).Class(`editq-fold-card`)
	if open { card = card.KV(`open`, `open`) }
	return card.Wrap(
		Elem(`summary`).Class(`editq-card-title`, `editq-fold-title`).Wrap(
			Div().Class(`editq-card-title-row`).Wrap(titleRow),
		),
		Div().Class(`editq-card-body`).Wrap(body),
	)
}

func EditQReviewControlValue(vars UIBagVars_t, name string) string {
	value := vars[name]
	if value == `` { return `-` }
	ctrl, ok := QuoteControlByName(name)
	if !ok { return value }
	if ctrl.kind == quoteCheckbox {
		if QuoteVarBool(value) { return `Yes` }
		return `No`
	}
	if ctrl.kind != quoteSelect { return value }
	for _, x := range QuoteControlChoices(ctrl) {
		if Str(x.id) == value { return x.label }
	}
	return value
}

func EditQReviewYear(value string) string {
	value = Trim(value)
	if len(value) < 4 { return `-` }
	return value[:4]
}

func EditQReviewClientView(vars UIBagVars_t) Elem_t {
	clientName := Trim(vars[`clientName`])
	if clientName == `` { clientName = `No client name` }
	return Div().Class(`editq-review-client`).Wrap(
		Div(clientName).Class(`editq-review-client-name`),
		Div().Class(`editq-review-client-grid`).Wrap(
			Div(`Year Born`).Class(`editq-review-client-key`),
			Div(EditQReviewYear(vars[`birth`])).Class(`editq-review-client-val`),
			Div(`Year Bought`).Class(`editq-review-client-key`),
			Div(EditQReviewYear(vars[`buy`])).Class(`editq-review-client-val`),
			Div(`Sick Cover`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `sickCover`)).Class(`editq-review-client-val`),
			Div(`Client Type`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `segment`)).Class(`editq-review-client-val`),
			Div(`Vision`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `vision`)).Class(`editq-review-client-val`),
		),
	)
}

func EditQReviewPreexByItemCateg(vars UIBagVars_t) (map[int]EuroCent_t, map[string]EuroCent_t) {
	byItem := make(map[int]EuroCent_t)
	byItemCateg := make(map[string]EuroCent_t)
	for _, x := range EditQPreexCharges(vars) {
		key := Str(x.itemId, `:`, x.categId)
		byItem[x.itemId] += x.applied
		byItemCateg[key] += x.applied
	}
	return byItem, byItemCateg
}

func EditQReviewPriceRow(label string, total, base, surch, preex EuroCent_t, class ...string) Elem_t {
	return Div().Class(`editq-review-row`).Class(class...).Wrap(
		Div(label).Class(`editq-review-col-label`),
		Div(EditQEuroCentText(total)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(base)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(surch)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(preex)).Class(`editq-review-col-money`, `editq-review-col-preex`),
	)
}

func EditQReviewPlanView(vars UIBagVars_t, item QuoteSelectedItem_t, row QuotePlan_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t) Elem_t {
	preexTotal := preexByItem[item.itemId]
	preexPlan := preexByItemCateg[Str(item.itemId, `:`, 0)]
	effectiveAge := `-`
	ageLabel := `Age`
	work := QuoteStateFromVars(vars)
	_, yearAge, exactAge := PlanAges(work)
	age := yearAge
	if EditQPlanUsesExactAge(row.planId) {
		age = exactAge
		ageLabel = `Exact age`
	}
	if age > 0 { effectiveAge = Str(age) }
	var rows []Elem_t
	rows = append(rows, Div().Class(`editq-review-row`, `editq-review-row-head`).Wrap(
		Div(`Item`).Class(`editq-review-col-label`),
		Div(`Total`).Class(`editq-review-col-money`),
		Div(`Base`).Class(`editq-review-col-money`),
		Div(`Surch`).Class(`editq-review-col-money`),
		Div(`Pre-ex`).Class(`editq-review-col-money`, `editq-review-col-preex`),
	))
	rows = append(rows, EditQReviewPriceRow(`Total`, row.base+row.surcharge+preexTotal, row.base, row.surcharge, preexTotal, `editq-review-row-total`))
	rows = append(rows, EditQReviewPriceRow(`Plan`, row.planBase+row.planSurcharge+preexPlan, row.planBase, row.planSurcharge, preexPlan))
	for _, addon := range row.addons {
		if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` { continue }
		label := QuoteAddonPickText(addon)
		if label == `` { label = addon.categ }
		preex := preexByItemCateg[Str(item.itemId, `:`, addon.categId)]
		rows = append(rows, EditQReviewPriceRow(label, addon.base+addon.surcharge+preex, addon.base, addon.surcharge, preex))
	}

	return Div().Class(`editq-review-plan`).Wrap(
		Div(row.label).Class(`editq-review-plan-title`),
		Div().Class(`editq-review-plan-meta`).Wrap(
			Span(Str(ageLabel, ` `, effectiveAge)),
			Span(Str(`Ded `, PriceTextWholeEuro(row.deduct, true))),
			Span(Str(`NC `, PriceTextWholeEuro(row.noClaims, true))),
			Span(Str(`Comm `, PriceTextWholeEuro(row.commission, true))),
		),
		Div().Class(`editq-review-grid`).Wrap(rows),
	)
}

func EditQReviewExportButtonsView() Elem_t {
	return Div().
		Class(`editq-card-title-right`).
		KV(`style`, `width: 100%; justify-content: center;`).
		Wrap(
			Elem(`button`).
				Type(`submit`).
				Name(`DownloadExcel`).
				Value(`slim=false`).
				KV(`formaction`, `/download-excel`).
				KV(`formmethod`, `post`).
				Class(`editq-title-btn`).
				Text(`Get Excel`),
			Elem(`button`).
				Type(`submit`).
				Name(`DownloadExcel`).
				Value(`slim=true`).
				KV(`formaction`, `/download-excel`).
				KV(`formmethod`, `post`).
				Class(`editq-title-btn`).
				Text(`Get Slim`),
		)
}

func EditQQuoteReviewBody(vars UIBagVars_t) Elem_t {
	state := QuoteStateFromVars(vars)
	selected := QuoteSelectedRows(state)
	preexByItem, preexByItemCateg := EditQReviewPreexByItemCateg(vars)

	var plans []Elem_t
	for _, x := range selected {
		plans = append(plans, EditQReviewPlanView(vars, x.item, x.row, preexByItem, preexByItemCateg))
	}
	if len(plans) == 0 {
		plans = append(plans, Div(`No plans selected.`).Class(`editq-review-empty`))
	}
	return Div().Class(`editq-review`).Wrap(
		EditQReviewExportButtonsView(),
		EditQReviewClientView(vars),
		Div().Class(`editq-review-plans`).Wrap(plans),
	)
}

func EditQHeaderView(vars UIBagVars_t) Elem_t {
	return EditQTopCardView(`EditQPreexCard`, EditQPreexTitleText(vars), false,
		Div().Class(`editq-header`).Wrap(
			EditQPreexChargesView(vars),
		),
	)
}

func EditQQuoteReviewCardView(vars UIBagVars_t) Elem_t {
	backBtn := Elem(`a`).
		KV(`href`, `/quote`).
		Class(`editq-title-btn`).
		Text(`Back to Quote Info`)
	return EditQTopCardView(`EditQReviewCard`, `Quote Review`, true,
		EditQQuoteReviewBody(vars),
		backBtn,
	)
}
