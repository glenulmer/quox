package main

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func QuotePlanAddonByTag(x QuotePlan_t, tag string) (QuotePlanAddon_t, bool) {
	tag = Lower(Trim(tag))
	for _, addon := range x.addons {
		if Contains(Lower(Trim(addon.categ)), tag) { return addon, true }
	}
	return QuotePlanAddon_t{}, false
}

func QuotePlanAddonByCateg(x QuotePlan_t, categId CategId_t) (QuotePlanAddon_t, bool) {
	for _, addon := range x.addons {
		if addon.categId == categId { return addon, true }
	}
	return QuotePlanAddon_t{}, false
}

func QuotePlanCellText(x QuotePlan_t, tag string) string {
	addon, ok := QuotePlanAddonByTag(x, tag)
	if !ok { return `` }
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		return ``
	}
	pick := QuoteAddonPickText(addon)
	price := PriceText(addon.base+addon.surcharge, addon.priceOk)
	switch {
	case pick != `` && price != `-`:
		return Str(pick, ` / `, price)
	case pick != ``:
		return pick
	default:
		return price
	}
}

func QuotePlanCategCellText(x QuotePlan_t, categId CategId_t) string {
	addon, ok := QuotePlanAddonByCateg(x, categId)
	if !ok { return `` }
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		return ``
	}
	pick := QuoteAddonPickText(addon)
	price := PriceText(addon.base+addon.surcharge, addon.priceOk)
	switch {
	case pick != `` && price != `-`:
		return Str(pick, ` / `, price)
	case pick != ``:
		return pick
	default:
		return price
	}
}

func QuoteAddonPickText(x QuotePlanAddon_t) string {
	if x.label != `` && Lower(Trim(x.label)) != Lower(Trim(x.categ)) { return x.label }
	if x.addon != 0 { return AddonName(CatChoice_t{ addon:x.addon, label:x.label }) }
	if x.level != 0 { return LevelLabel(x.level) }
	return ``
}

func QuotePlanAddonSelectView(planId int, x QuotePlanAddon_t) Elem_t {
	key := QuotePlanCatControlName(planId, x.categId)
	var options []Elem_t
	for _, choice := range x.choices {
		options = append(options, Option().KV(`value`, choice.addon).Text(AddonName(choice)))
	}
	return Select(options).Id(key).Name(key).Choose(x.addon)
}

func QuotePlanAddonView(planId int, x QuotePlanAddon_t) Elem_t {
	pick := QuoteAddonPickText(x)
	pickView := any(pick)
	if x.hasMulti && len(x.choices) > 0 {
		pickView = QuotePlanAddonSelectView(planId, x)
	}
	return Div().Class(`quote-plan-addon`).Wrap(
		Div(x.categ).Class(`quote-plan-addon-cat`),
		Div().Class(`quote-plan-addon-pick`).Wrap(pickView),
		Div(PriceText(x.base, x.priceOk)).Class(`quote-plan-addon-base`),
		Div(PriceText(x.surcharge, x.priceOk)).Class(`quote-plan-addon-surch`),
	)
}

func QuotePlanBaseView(x QuotePlan_t) Elem_t {
	return Div().Class(`quote-plan-addon`, `quote-plan-addon-base-row`).Wrap(
		Div(`Plan`).Class(`quote-plan-addon-cat`),
		Div(``).Class(`quote-plan-addon-pick`),
		Div(PriceText(x.planBase, x.planOk)).Class(`quote-plan-addon-base`),
		Div(PriceText(x.planSurcharge, x.planOk)).Class(`quote-plan-addon-surch`),
	)
}

func QuotePlanCardView(x QuotePlan_t) Elem_t {
	var addons []Elem_t
	addons = append(addons, QuotePlanBaseView(x))
	for _, addon := range x.addons { addons = append(addons, QuotePlanAddonView(x.planId, addon)) }
	return Div().Class(`quote-plan-card`).Wrap(
		Div().Class(`quote-plan-head`).Wrap(
			Div(x.label).Class(`quote-plan-label`),
			Div().Class(`quote-plan-total`).Wrap(
				Div(`Total`).Class(`quote-plan-total-label`),
				Div(PriceText(x.price, true)).Class(`quote-plan-price`),
			),
		),
		Div().Class(`quote-plan-sums`).Wrap(
			Div(`Base: `, PriceText(x.base, true)),
			Div(`Surcharge: `, PriceText(x.surcharge, true)),
		),
		Div().Class(`quote-plan-addon-list`).Wrap(addons),
	)
}

func QuoteFilteredPlanView(x QuotePlanFiltered_t) Elem_t {
	return Div().Class(`quote-filtered-plan`).Wrap(
		Div(x.label).Class(`quote-filtered-label`),
		Div(Join(x.reasons, `, `)).Class(`quote-filtered-reasons`),
	)
}

func QuotePlanDesktopCategs() []Categ_t {
	var out []Categ_t
	for _, categ := range App.lookup.categs.All() {
		if categ.categId == 0 { continue }
		out = append(out, categ)
	}
	return out
}

func QuotePlanDesktopGridStyle(categCount int) string {
	x := `grid-template-columns: 70px 255px 92px`
	for i := 0; i < categCount; i++ {
		x += ` 120px`
	}
	x += ` 90px;`
	return x
}

func QuotePlanDesktopHead(categs []Categ_t) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div(`Total`).Class(`quote-plan-cell`, `quote-plan-total-cell`),
		Div(`Plan`).Class(`quote-plan-cell`, `quote-plan-name-cell`),
		Div(`Plan row`).Class(`quote-plan-cell`),
	)
	for _, categ := range categs {
		cols = append(cols, Div(categ.name).Class(`quote-plan-cell`))
	}
	cols = append(cols, Div(`Vision`).Class(`quote-plan-cell`))
	return Div().
		Class(`quote-plan-table-row`, `quote-plan-table-head`).
		KV(`style`, QuotePlanDesktopGridStyle(len(categs))).
		Wrap(cols)
}

func QuotePlanDesktopRow(x QuotePlan_t, categs []Categ_t) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div(PriceText(x.price, true)).Class(`quote-plan-cell`, `quote-plan-total-cell`),
		Div(x.label).Class(`quote-plan-cell`, `quote-plan-name-cell`),
		Div(PriceText(x.planBase+x.planSurcharge, x.planOk)).Class(`quote-plan-cell`),
	)
	for _, categ := range categs {
		cols = append(cols, Div(QuotePlanCategCellText(x, categ.categId)).Class(`quote-plan-cell`))
	}
	cols = append(cols, Div(QuotePlanCellText(x, `vision`)).Class(`quote-plan-cell`))
	return Div().
		Class(`quote-plan-table-row`).
		KV(`style`, QuotePlanDesktopGridStyle(len(categs))).
		Wrap(cols)
}

func QuotePlanDesktopView(data QuotePlans_t) Elem_t {
	categs := QuotePlanDesktopCategs()
	var rows []Elem_t
	rows = append(rows, QuotePlanDesktopHead(categs))
	for _, x := range data.plans { rows = append(rows, QuotePlanDesktopRow(x, categs)) }
	return Div().Class(`quote-plan-table`).Wrap(rows)
}

func QuoteFilteredPlansBox(filteredPlans []QuotePlanFiltered_t) Elem_t {
	var filtered []Elem_t
	for _, x := range filteredPlans { filtered = append(filtered, QuoteFilteredPlanView(x)) }
	if len(filtered) == 0 {
		filtered = append(filtered, Div(`No plans filtered out.`).Class(`quote-filtered-empty`))
	}
	return Elem(`details`).Class(`quote-filter-box`).Wrap(
		Elem(`summary`).Class(`quote-filter-box-title`).Wrap(
			Span(`Filtered plans (` , len(filteredPlans), `)`),
		),
		Div().Class(`quote-filtered-list`).Wrap(filtered),
	)
}
