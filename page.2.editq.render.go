package main

import (
	. "pm/lib/dec2"
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
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

func EditQHeaderView(vars QuoteVars_t) Elem_t {
	custName := vars[`custName`]
	return EditQTopCardView(`EditQPrexCard`, EditQPrimeTitleText(vars), false,
		Div().Class(`editq-header`).Wrap(
			Elem(`label`).Class(`editq-field`, `editq-customer-field`).Wrap(
				QuoteInputText(`custName`, custName, `Customer name`),
			),
			EditQPrimeChargesView(vars),
		),
	)
}

func EditQQuoteReviewCardView() Elem_t {
	backBtn := Elem(`a`).
		KV(`href`, `/quote`).
		Class(`editq-title-btn`).
		Text(`Back to Quote Info`)
	return EditQTopCardView(`EditQReviewCard`, `Quote Review`, true,
		Div().Class(`editq-quote-review-empty`),
		backBtn,
	)
}
