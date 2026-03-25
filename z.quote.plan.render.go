package main

import (
	"fmt"
	"strings"
	. "pm/lib/htmlHelper"
	. "pm/lib/dec2"
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
				Div(PriceTextWholeEuro(x.price, true)).Class(`quote-plan-price`),
			),
		),
		Div().Class(`quote-plan-meta`).Wrap(
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Ded`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.deduct, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`NC`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.noClaims, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Comm`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.commission, true)).Class(`quote-plan-meta-value`),
			),
		),
		Elem(`details`).Class(`quote-plan-addon-details`).Wrap(
			Elem(`summary`).Class(`quote-plan-addon-title`).Wrap(
				Span(`Addon Prices`).Class(`quote-plan-addon-title-label`),
				Span(PriceText(x.base, true)).Class(`quote-plan-addon-title-sum`),
				Span(PriceText(x.surcharge, true)).Class(`quote-plan-addon-title-sum`),
			),
			Div().Class(`quote-plan-addon-list`).Wrap(addons),
		),
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
		if categ.display == 0 { continue }
		out = append(out, categ)
	}
	return out
}

func QuoteSortSelectView(sortBy string) Elem_t {
	mode := QuoteSortMode(sortBy)
	return Select(
		Option().KV(`value`, sortByName).Text(`Name`),
		Option().KV(`value`, sortByPrice).Text(`Total`),
	).Name(`sortBy`).Choose(mode).Class(`quote-plan-sort-input`)
}

func QuotePlanDesktopCategWidth(categ Categ_t) int {
	name := Lower(Trim(categ.name))
	switch {
	case Contains(name, `sick`):
		return 72
	case Contains(name, `pvn`):
		return 50
	case Contains(name, `hospital`):
		return 108
	case Contains(name, `dental`):
		return 108
	case Contains(name, `consumer`):
		return 84
	case Contains(name, `natural`):
		return 80
	case Contains(name, `special`):
		return 90
	}
	return 120
}

func QuotePxToRem(px int) string {
	s := fmt.Sprintf(`%.4f`, float64(px)/16.0)
	s = strings.TrimRight(s, `0`)
	s = strings.TrimRight(s, `.`)
	return s + `rem`
}

func QuotePlanDesktopGridStyle(categs []Categ_t, showVision bool) string {
	const moneyColPx = 67
	x := Str(`grid-template-columns: `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(294))
	for _, categ := range categs {
		x += Str(` `, QuotePxToRem(QuotePlanDesktopCategWidth(categ)))
	}
	if showVision {
		x += Str(` `, QuotePxToRem(60))
	}
	x += Str(` `, QuotePxToRem(moneyColPx))
	x += `;`
	return x
}

func QuotePlanDesktopHead(categs []Categ_t, showVision bool, sortBy string, planCount int) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div(`Total`).Class(`quote-plan-cell`, `quote-plan-total-cell`),
		Div(`Ded`).Class(`quote-plan-cell`),
		Div(`NC`).Class(`quote-plan-cell`, `quote-plan-money-head-right`),
		Div().Class(`quote-plan-cell`, `quote-plan-name-cell`).Wrap(
			Div().Class(`quote-plan-head-plan`).Wrap(
				Span(`Plans (` , planCount , `)`).Class(`quote-plan-head-plan-title`),
				QuoteSortSelectView(sortBy),
			),
		),
	)
	for _, categ := range categs {
		cols = append(cols, Div(categ.name).Class(`quote-plan-cell`))
	}
	if showVision {
		cols = append(cols, Div(`Vision`).Class(`quote-plan-cell`))
	}
	cols = append(cols, Div(`Comm`).Class(`quote-plan-cell`, `quote-plan-money-head-right`))
	return Div().
		Class(`quote-plan-table-row`, `quote-plan-table-head`).
		KV(`style`, QuotePlanDesktopGridStyle(categs, showVision)).
		Wrap(cols)
}

func QuotePlanDesktopAddonPickView(planId int, addon QuotePlanAddon_t) Elem_t {
	if addon.hasMulti && len(addon.choices) > 0 {
		return Div().Class(`quote-plan-cell-pick`).Wrap(QuotePlanAddonSelectView(planId, addon))
	}
	return Div(QuoteAddonPickText(addon)).Class(`quote-plan-cell-pick`)
}

func QuotePlanDesktopCategCellView(x QuotePlan_t, categId CategId_t) Elem_t {
	addon, ok := QuotePlanAddonByCateg(x, categId)
	if !ok { return Div().Class(`quote-plan-cell-pick`) }
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		return Div().Class(`quote-plan-cell-pick`)
	}
	return QuotePlanDesktopAddonPickView(x.planId, addon)
}

func QuotePlanDesktopTagCellView(x QuotePlan_t, tag string) Elem_t {
	addon, ok := QuotePlanAddonByTag(x, tag)
	if !ok { return Div().Class(`quote-plan-cell-pick`) }
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		return Div().Class(`quote-plan-cell-pick`)
	}
	return QuotePlanDesktopAddonPickView(x.planId, addon)
}

func QuotePlanDesktopVisionCellView(x QuotePlan_t) Elem_t {
	addon, ok := QuotePlanAddonByTag(x, `vision`)
	if !ok { return Div().Class(`quote-plan-cell-pick`) }
	if !addon.priceOk { return Div().Class(`quote-plan-cell-pick`) }
	return Div(PriceText(addon.base+addon.surcharge, addon.priceOk)).Class(`quote-plan-cell-pick`, `quote-plan-cell-money`)
}

func PriceTextWholeEuro(amount EuroCent_t, ok bool) string {
	if !ok { return `-` }
	return EuroFlatFromCent(amount).OutEuro()
}

func QuotePlanDesktopRow(x QuotePlan_t, categs []Categ_t, showVision bool) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div(PriceTextWholeEuro(x.price, true)).Class(`quote-plan-cell`, `quote-plan-total-cell`),
		Div(PriceTextWholeEuro(x.deduct, true)).Class(`quote-plan-cell`, `quote-plan-cell-money`),
		Div(PriceTextWholeEuro(x.noClaims, true)).Class(`quote-plan-cell`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(x.label).Class(`quote-plan-cell`, `quote-plan-name-cell`),
	)
	for _, categ := range categs {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopCategCellView(x, categ.categId)))
	}
	if showVision {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopVisionCellView(x)))
	}
	cols = append(cols, Div(PriceTextWholeEuro(x.commission, true)).Class(`quote-plan-cell`, `quote-plan-cell-money`, `quote-plan-cell-money-right`))
	return Div().
		Class(`quote-plan-table-row`).
		KV(`style`, QuotePlanDesktopGridStyle(categs, showVision)).
		Wrap(cols)
}

func QuotePlanDesktopView(data QuotePlans_t) Elem_t {
	categs := QuotePlanDesktopCategs()
	var rows []Elem_t
	rows = append(rows, QuotePlanDesktopHead(categs, data.showVision, data.sortBy, len(data.plans)))
	for _, x := range data.plans { rows = append(rows, QuotePlanDesktopRow(x, categs, data.showVision)) }
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
