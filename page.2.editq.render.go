package main

import (
	. "quo2/lib/dec2"
	. "quo2/lib/htmlHelper"
	. "quo2/lib/output"
)

func EditQGrowInput(name, value string) Elem_t {
	return Elem(`textarea`).
		Name(name).
		Class(`editq-grow-input`).
		KV(`rows`, 1).
		Text(value)
}

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

func EditQConditionRow(name, value, delName string) Elem_t {
	return Div().Class(`editq-condition-row`).Wrap(
		Div().Class(`editq-condition-text`).Wrap(EditQGrowInput(name, value)),
		Div().Class(`editq-condition-del`).Wrap(EditQDelButton(delName)),
	)
}

func EditQPrimeChargeModeView(name, mode string) Elem_t {
	return Select(
		Option().KV(`value`, editQPrimeModePct).Text(`%`),
		Option().KV(`value`, editQPrimeModeEur).Text(`€`),
	).Name(name).Choose(EditQPrimeMode(mode)).Class(`editq-prime-mode`)
}

func EditQPrimeChargeRowView(x EditQPrimeCharge_t, modeKey, amountKey, noteKey string) Elem_t {
	return Div().Class(`editq-prime-row`).Wrap(
		Div(x.level).Class(`editq-prime-categ`),
		Div().Class(`editq-prime-inputs`).Wrap(
			EditQPrimeChargeModeView(modeKey, x.mode),
			QuoteInputText(amountKey, x.amount, `Amount`).Class(`editq-prime-amount`, `editq-prime-amount-input`),
			QuoteInputText(noteKey, x.note, `Optional note`).Class(`editq-prime-note`),
		),
	)
}

func EditQEuroCentText(v EuroCent_t) string {
	return Str(v.String(), ` €`)
}

func EditQPrimeTotal(vars QuoteVars_t) EuroCent_t {
	charges := EditQPrimeCharges(vars)
	var total EuroCent_t
	for _, x := range charges { total += x.applied }
	return total
}

func EditQPrimeTitleText(vars QuoteVars_t) string {
	totalFlat := EuroFlatFromCent(EditQPrimeTotal(vars))
	return Str(`Pre-existing conditions (`, totalFlat.OutEuro(), `)`)
}

func EditQPrimeChargesView(vars QuoteVars_t) Elem_t {
	charges := EditQPrimeCharges(vars)
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
				prex := appliedByItem[prevItemId]
				rows = append(rows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-prime-summary-base`),
					Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
						Div().Class(`editq-prime-summary-spacer`),
						Div(EditQEuroCentText(prex)).Class(`editq-prime-summary-prex`),
						Div(Str(`= `, EditQEuroCentText(base+prex))).Class(`editq-prime-summary-sum`),
					),
				))
			}
			rows = append(rows, Div().Class(`editq-prime-plan`).Wrap(
				Span(x.plan).Class(`editq-prime-plan-name`),
			))
			prevItemId = x.itemId
			hasPlan = true
		}
		rows = append(rows, EditQPrimeChargeRowView(x, EditQPrimeModeKey(x.itemId, x.categId), EditQPrimeAmountKey(x.itemId, x.categId), EditQPrimeNoteKey(x.itemId, x.categId)))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		prex := appliedByItem[prevItemId]
		rows = append(rows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-prime-summary-base`),
			Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
				Div().Class(`editq-prime-summary-spacer`),
				Div(EditQEuroCentText(prex)).Class(`editq-prime-summary-prex`),
				Div(Str(`= `, EditQEuroCentText(base+prex))).Class(`editq-prime-summary-sum`),
			),
		))
	}
	if len(rows) == 0 {
		rows = append(rows, Div(`No payable selected plan categories yet.`).Class(`editq-prime-empty`))
	}
	return Div().Class(`editq-section`, `editq-prime-charges`).Wrap(
		Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
		Div().Class(`editq-prime-list`).Wrap(rows),
	)
}

func EditQDependentView(vars QuoteVars_t, dep EditQDep_t, order int) Elem_t {
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
				prex := appliedByItem[prevItemId]
				depRows = append(depRows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-prime-summary-base`),
					Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
						Div().Class(`editq-prime-summary-spacer`),
						Div(EditQEuroCentText(prex)).Class(`editq-prime-summary-prex`),
						Div(Str(`= `, EditQEuroCentText(base+prex))).Class(`editq-prime-summary-sum`),
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
				depRows = append(depRows, Div().Class(`editq-prime-plan`).Wrap(
					Span(x.plan).Class(`editq-prime-plan-name`),
				))
			} else {
				depRows = append(depRows, Div().Class(`editq-prime-plan`).Wrap(
					Div().Class(`editq-prime-plan-head`).Wrap(
						Span(x.plan).Class(`editq-prime-plan-name`),
					),
					Div(ageText).Class(`editq-prime-plan-age`),
				))
			}
			prevItemId = x.itemId
			hasPlan = true
		}
		depRows = append(depRows, EditQPrimeChargeRowView(
			x,
			EditQDepChargeModeKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeAmountKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeNoteKey(dep.depId, x.itemId, x.categId),
		))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		prex := appliedByItem[prevItemId]
		depRows = append(depRows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-prime-summary-base`),
			Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
				Div().Class(`editq-prime-summary-spacer`),
				Div(EditQEuroCentText(prex)).Class(`editq-prime-summary-prex`),
				Div(Str(`= `, EditQEuroCentText(base+prex))).Class(`editq-prime-summary-sum`),
			),
		))
	}
	if len(depRows) == 0 {
		depRows = append(depRows, Div(`No selected plan available.`).Class(`editq-prime-empty`))
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
		Div().Class(`editq-prime-charges`, `editq-dependent-prex`).Wrap(
			Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
			Div().Class(`editq-prime-list`).Wrap(depRows),
		),
	)
}

func EditQDependentsView(vars QuoteVars_t, sortForGet bool) Elem_t {
	deps := EditQDependents(vars, sortForGet)
	namedCount := 0
	var list []Elem_t
	for i, dep := range deps {
		if Trim(dep.name) != `` { namedCount++ }
		list = append(list, EditQDependentView(vars, dep, i))
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

func EditQReviewControlValue(vars QuoteVars_t, name string) string {
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

func EditQReviewClientView(vars QuoteVars_t) Elem_t {
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

func EditQReviewPrexByItemCateg(vars QuoteVars_t) (map[int]EuroCent_t, map[string]EuroCent_t) {
	byItem := make(map[int]EuroCent_t)
	byItemCateg := make(map[string]EuroCent_t)
	for _, x := range EditQPrimeCharges(vars) {
		key := Str(x.itemId, `:`, x.categId)
		byItem[x.itemId] += x.applied
		byItemCateg[key] += x.applied
	}
	return byItem, byItemCateg
}

func EditQReviewPriceRow(label string, total, base, surch, prex EuroCent_t, class ...string) Elem_t {
	return Div().Class(`editq-review-row`).Class(class...).Wrap(
		Div(label).Class(`editq-review-col-label`),
		Div(EditQEuroCentText(total)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(base)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(surch)).Class(`editq-review-col-money`),
		Div(EditQEuroCentText(prex)).Class(`editq-review-col-money`, `editq-review-col-prex`),
	)
}

func EditQReviewPlanView(vars QuoteVars_t, item QuoteSelectedItem_t, row QuotePlan_t, prexByItem map[int]EuroCent_t, prexByItemCateg map[string]EuroCent_t) Elem_t {
	prexTotal := prexByItem[item.itemId]
	prexPlan := prexByItemCateg[Str(item.itemId, `:`, 0)]
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
		Div(`Pre-ex`).Class(`editq-review-col-money`, `editq-review-col-prex`),
	))
	rows = append(rows, EditQReviewPriceRow(`Total`, row.base+row.surcharge+prexTotal, row.base, row.surcharge, prexTotal, `editq-review-row-total`))
	rows = append(rows, EditQReviewPriceRow(`Plan`, row.planBase+row.planSurcharge+prexPlan, row.planBase, row.planSurcharge, prexPlan))
	for _, addon := range row.addons {
		if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` { continue }
		label := QuoteAddonPickText(addon)
		if label == `` { label = addon.categ }
		prex := prexByItemCateg[Str(item.itemId, `:`, addon.categId)]
		rows = append(rows, EditQReviewPriceRow(label, addon.base+addon.surcharge+prex, addon.base, addon.surcharge, prex))
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
			Elem(`button`).Type(`button`).Class(`editq-title-btn`).Text(`Get Excel`),
			Elem(`button`).Type(`button`).Class(`editq-title-btn`).Text(`Get Slim`),
		)
}

func EditQQuoteReviewBody(vars QuoteVars_t) Elem_t {
	state := QuoteStateFromVars(vars)
	selected := QuoteSelectedRows(state)
	prexByItem, prexByItemCateg := EditQReviewPrexByItemCateg(vars)

	var plans []Elem_t
	for _, x := range selected {
		plans = append(plans, EditQReviewPlanView(vars, x.item, x.row, prexByItem, prexByItemCateg))
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

func EditQHeaderView(vars QuoteVars_t) Elem_t {
	return EditQTopCardView(`EditQPrexCard`, EditQPrimeTitleText(vars), false,
		Div().Class(`editq-header`).Wrap(
			EditQPrimeChargesView(vars),
		),
	)
}

func EditQQuoteReviewCardView(vars QuoteVars_t) Elem_t {
	backBtn := Elem(`a`).
		KV(`href`, `/quote`).
		Class(`editq-title-btn`).
		Text(`Back to Quote Info`)
	return EditQTopCardView(`EditQReviewCard`, `Quote Review`, true,
		EditQQuoteReviewBody(vars),
		backBtn,
	)
}
