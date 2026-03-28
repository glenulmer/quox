package main

import (
	"fmt"

	. "pm/lib/htmlHelper"
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

func EditQRootConditionsView(vars QuoteVars_t) Elem_t {
	conditions := EditQPreConditions(vars)
	var rows []Elem_t
	for _, x := range conditions {
		rows = append(rows, EditQConditionRow(EditQPreKey(x.condId), x.text, EditQPreDelControlName(x.condId)))
	}
	rows = append(rows, EditQAddButton(EditQPreAddControlName(), `Pre-existing`))
	return Div().Class(`editq-section`, `editq-preexisting`).Wrap(
		Div(`Pre-existing Conditions`).Class(`editq-section-title`),
		Div().Class(`editq-condition-list`).Wrap(rows),
	)
}

func EditQDependentView(dep EditQDep_t, order int) Elem_t {
	var condRows []Elem_t
	for _, x := range dep.conds {
		condRows = append(condRows, EditQConditionRow(EditQDepCondKey(dep.depId, x.condId), x.text, EditQDepPreDelControlName(dep.depId, x.condId)))
	}
	condRows = append(condRows, EditQAddButton(EditQDepPreAddControlName(dep.depId), `Pre-existing`, `editq-dependent-pre-btn`))

	return Div().Class(`editq-dependent`).Wrap(
		Div().Class(`editq-dependent-head`).Wrap(
			Div(fmt.Sprintf(`Dependent %d`, order+1)).Class(`editq-dependent-title`),
			EditQDelButton(EditQDepDelControlName(dep.depId)),
		),
		Div().Class(`editq-dependent-fields`).Wrap(
			Elem(`label`).Class(`editq-field`).Wrap(
				Span(`Name`).Class(`editq-label`),
				QuoteInputText(EditQDepNameKey(dep.depId), dep.name, `Dependent name`),
			),
			Elem(`label`).Class(`editq-field`).Wrap(
				Span(`Birth date`).Class(`editq-label`),
				QuoteInputDate(EditQDepBirthKey(dep.depId), dep.birth),
			),
			Elem(`label`).Class(`editq-check`).Wrap(
				QuoteCheckbox(EditQDepVisionKey(dep.depId), dep.vision),
				Span(`Vision`).Class(`editq-check-text`),
			),
		),
		Div().Class(`editq-condition-list`, `editq-dependent-conditions`).Wrap(condRows),
	)
}

func EditQDependentsView(vars QuoteVars_t, sortForGet bool) Elem_t {
	deps := EditQDependents(vars, sortForGet)
	var list []Elem_t
	for i, dep := range deps {
		list = append(list, EditQDependentView(dep, i))
	}
	if len(deps) < editQDepMaxCount {
		list = append(list, EditQAddButton(EditQDepAddControlName(), `Add Dependent`, `editq-add-dependent`))
	}
	return Div().Class(`editq-section`, `editq-dependents`).Wrap(
		Div(`Dependents`).Class(`editq-section-title`),
		Div().Class(`editq-dependent-list`).Wrap(list),
	)
}

func EditQHeaderView(vars QuoteVars_t) Elem_t {
	custName := vars[`custName`]
	preview := Div(`No customer name set.`).Class(`editq-cust-preview`, `editq-cust-empty`)
	if custName != `` {
		preview = Div(`Customer: `, custName).Class(`editq-cust-preview`)
	}
	return Div().Class(`editq-header`).Wrap(
		Div().Class(`editq-header-top`).Wrap(
			Div(`Edit Quote`).Class(`editq-title`),
			Elem(`a`).KV(`href`, `/quote`).Class(`editq-back-link`).Text(`Back to Quote`),
		),
		preview,
		Elem(`label`).Class(`editq-field`, `editq-customer-field`).Wrap(
			Span(`Customer name`).Class(`editq-label`),
			QuoteInputText(`custName`, custName, `Customer name`),
		),
	)
}
