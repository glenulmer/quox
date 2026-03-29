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
				prexText := prex.OutEuro()
				if prex == 0 { prexText = `0,00 €` }
				rows = append(rows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
					Div(base.OutEuro()).Class(`editq-prime-summary-base`),
					Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
						Div().Class(`editq-prime-summary-spacer`),
						Div(prexText).Class(`editq-prime-summary-prex`),
						Div(Str(`= `, (base+prex).OutEuro())).Class(`editq-prime-summary-sum`),
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
		prexText := prex.OutEuro()
		if prex == 0 { prexText = `0,00 €` }
		rows = append(rows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
			Div(base.OutEuro()).Class(`editq-prime-summary-base`),
			Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
				Div().Class(`editq-prime-summary-spacer`),
				Div(prexText).Class(`editq-prime-summary-prex`),
				Div(Str(`= `, (base+prex).OutEuro())).Class(`editq-prime-summary-sum`),
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
				prexText := prex.OutEuro()
				if prex == 0 { prexText = `0,00 €` }
				depRows = append(depRows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
					Div(base.OutEuro()).Class(`editq-prime-summary-base`),
					Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
						Div().Class(`editq-prime-summary-spacer`),
						Div(prexText).Class(`editq-prime-summary-prex`),
						Div(Str(`= `, (base+prex).OutEuro())).Class(`editq-prime-summary-sum`),
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
		prexText := prex.OutEuro()
		if prex == 0 { prexText = `0,00 €` }
		depRows = append(depRows, Div().Class(`editq-prime-row`, `editq-prime-summary`).Wrap(
			Div(base.OutEuro()).Class(`editq-prime-summary-base`),
			Div().Class(`editq-prime-inputs`, `editq-prime-summary-inputs`).Wrap(
				Div().Class(`editq-prime-summary-spacer`),
				Div(prexText).Class(`editq-prime-summary-prex`),
				Div(Str(`= `, (base+prex).OutEuro())).Class(`editq-prime-summary-sum`),
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
	var list []Elem_t
	for i, dep := range deps {
		list = append(list, EditQDependentView(vars, dep, i))
	}
	if len(deps) < editQDepMaxCount {
		list = append(list, EditQAddButton(EditQDepAddControlName(), `Add`, `editq-add-dependent`))
	}
	return Div().Class(`editq-section`, `editq-dependents`).Wrap(
		Div(`Dependents`).Class(`editq-section-title`),
		Div().Class(`editq-dependent-list`).Wrap(list),
	)
}

func EditQHeaderView(vars QuoteVars_t) Elem_t {
	custName := vars[`custName`]
	return Div().Class(`editq-header`).Wrap(
		Div().Class(`editq-header-top`).Wrap(
			Div(`Quote Review`).Class(`editq-title`),
			Elem(`a`).KV(`href`, `/quote`).Class(`editq-back-link`).Text(`Back to Quote Info`),
		),
		Elem(`label`).Class(`editq-field`, `editq-customer-field`).Wrap(
			QuoteInputText(`custName`, custName, `Customer name`),
		),
		EditQPrimeChargesView(vars),
	)
}
