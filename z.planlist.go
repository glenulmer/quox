package main

import (
	. "pm/lib/date"
	. "pm/lib/dec2"
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

const catSick CategId_t = 1

func PlanCard(planId int, plan Plan_t, priceText, sumBaseText, sumSurchText string, addonRows []Elem_t) Elem_t {
	body := Div().Class(`card-body`, `plan-card-body`).Wrap(
		Div().Class(`plan-line`).Wrap(
			Div(plan.provName, ` / `, plan.name).Class(`plan-label`),
			Div(priceText).Class(`plan-price`),
		),
	)
	if len(addonRows) > 0 {
		body = body.Wrap(
			Elem(`details`).Class(`plan-addon-details`).Wrap(
				Elem(`summary`).Class(`plan-addon-title`).Wrap(
					Span(`Total`).Class(`plan-addon-cat`, `plan-addon-title-label`),
					Span(``).Class(`plan-addon-pick`),
					Span(sumBaseText).Class(`plan-addon-base`, `plan-addon-title-sum`),
					Span(sumSurchText).Class(`plan-addon-surch`, `plan-addon-title-sum`),
				),
				Div().Class(`plan-addon-body`).Wrap(addonRows),
			),
		)
	}
	return Div().Class(`card`).Id(`planId`, planId).Wrap(body)
}

func PlanAddonHead() Elem_t {
	return Div().Class(`plan-addon-head`).Wrap(
		Div(``).Class(`plan-addon-cat`),
		Div(``).Class(`plan-addon-pick`),
		Div(`base`).Class(`plan-addon-base`),
		Div(`surch`).Class(`plan-addon-surch`),
	)
}

func PlanAddonRow(cat, pick any, baseText, surchText string, class ...string) Elem_t {
	row := Div().Class(`plan-addon-row`)
	if len(class) > 0 { row = row.Class(class...) }
	return row.Wrap(
		Div().Class(`plan-addon-cat`).Wrap(cat),
		Div().Class(`plan-addon-pick`).Wrap(pick),
		Div(baseText).Class(`plan-addon-base`),
		Div(surchText).Class(`plan-addon-surch`),
	)
}

func PlanCatSelect(planId int, categ Categ_t, choices []CatChoice_t, selected AddonId_t) Elem_t {
	var options []Elem_t
	for _, choice := range choices {
		options = append(options, Option().KV(`value`, choice.addon).Text(AddonName(choice)))
	}
	return Select(options).
		Id(`plancat-`, planId, `-`, categ.categId).
		SelO(selected)
}

func StateValue(state State_t, key string) string {
	v := state.quote[key]
	if v == `` { v = state.quote[Q(key)] }
	v = Trim(v)
	if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' { v = v[1:len(v)-1] }
	return v
}

func StateInt(state State_t, key string) int {
	return Atoi(StateValue(state, key))
}

func StateDate(state State_t, key string) CalDate_t {
	return Parse(`yyyy-mm-dd`, StateValue(state, key))
}

func PlanAges(state State_t) (buyYear, yearAge, exactAge int) {
	buy := StateDate(state, `buy`)
	birth := StateDate(state, `birth`)
	if !Valid(buy) || !Valid(birth) { return 0, 0, 0 }

	buyYear = buy.Year()
	yearAge = buyYear - birth.Year()
	exactAge = yearAge
	if int(buy)%10000 < int(birth)%10000 { exactAge-- }
	return buyYear, yearAge, exactAge
}

func LookupPrice(buyYear, age, productId int) (base, surcharge EuroCent_t, ok bool) {
	price, ok := App.lookup.prices[YAP(buyYear, age, productId)]
	if !ok { return 0, 0, false }
	return price.base, price.surcharge, true
}

func SickMult(sickCover int) EuroCent_t {
	return EuroCent_t(sickCover / 4500)
}

func CatPrice(buyYear, age, productId int, categId CategId_t, sickCover int) (base, surcharge EuroCent_t, ok bool) {
	base, surcharge, ok = LookupPrice(buyYear, age, productId)
	if !ok { return 0, 0, false }
	if categId == catSick { base *= SickMult(sickCover) }
	return base, surcharge, true
}

func PriceText(amount EuroCent_t, ok bool) string {
	if !ok { return `-` }
	return amount.String() + ` €`
}

func AddonName(choice CatChoice_t) string {
	label := Trim(choice.label)
	if label != `` { return label }
	product, ok := App.lookup.products[ProductId_t(choice.addon)]
	if ok { return product.name }
	return Str(choice.addon)
}

func LogFirstPlanAddonChoices(categ Categ_t, choices []CatChoice_t, buyYear, age, sickCover int) {
	for _, choice := range choices {
		base, surch, ok := CatPrice(buyYear, age, int(choice.addon), categ.categId, sickCover)
		star := ``
		if choice.isdef { star = ` *` }
		Log(`planList first plan addon:`,
			`categ=`, categ.name,
			`addon=`, AddonName(choice) + star,
			`base=`, PriceText(base, ok),
			`surcharge=`, PriceText(surch, ok),
		)
	}
}

func ListPlans(state State_t) Elem_t {
	buyYear, yearAge, exactAge := PlanAges(state)
	sickCover := StateInt(state, `sickCover`)

	var list []Elem_t
	loggedFirst := false
	for planId, plan := range App.lookup.plans.All() {
		age := yearAge
		if plan.exactAge { age = exactAge }
		firstPlan := !loggedFirst

		planBase, planSurch, planOk := CatPrice(buyYear, age, planId, 0, sickCover)
		sumBase, sumSurch := EuroCent_t(0), EuroCent_t(0)
		if planOk { sumBase += planBase; sumSurch += planSurch }

		addonRows := []Elem_t{
			PlanAddonHead(),
			PlanAddonRow(`Plan`, ``, PriceText(planBase, planOk), PriceText(planSurch, planOk)),
		}

		planKey := plan.planId
		for _, categ := range App.lookup.categs.All() {
			if categ.categId == 0 { continue }
			key := PlanCateg_t{ plan: planKey, categ: categ.categId }
			choices := App.lookup.planAddonChoices[key]
			if len(choices) == 0 { continue }

			if firstPlan { LogFirstPlanAddonChoices(categ, choices, buyYear, age, sickCover) }

			choice, ok := App.lookup.planAddons[key]
			if !ok { continue }
			hasMulti := len(choices) > 1
			if !hasMulti && choice.addon == 0 { continue }

			var base, surch EuroCent_t
			priceOk := false
			if choice.addon != 0 {
				base, surch, priceOk = CatPrice(buyYear, age, int(choice.addon), categ.categId, sickCover)
				if priceOk { sumBase += base; sumSurch += surch }
			}

			rowPick := any(``)
			label := Trim(choice.label)
			if label != `` && Lower(label) != Lower(Trim(categ.name)) { rowPick = label }
			if hasMulti {
				rowPick = PlanCatSelect(planId, categ, choices, choice.addon)
			}
			addonRows = append(addonRows,
				PlanAddonRow(categ.name, rowPick, PriceText(base, priceOk), PriceText(surch, priceOk)),
			)
		}

		sumBaseText := PriceText(sumBase, true)
		sumSurchText := PriceText(sumSurch, true)
		totalText := PriceText(sumBase + sumSurch, true)
		if firstPlan {
			Log(`planList first plan:`, `yearAge=`, yearAge, `planId=`, planId, `price=`, totalText)
		}
		list = append(list, PlanCard(planId, plan, totalText, sumBaseText, sumSurchText, addonRows))
		loggedFirst = true
	}

	return Div().Id(`planList`).Class(`card`).Wrap(
		Div(`Plans`).Class(`card-title`),
		Div().Class(`card-body`, `plan-list-body`).Wrap(list),
	)
}
