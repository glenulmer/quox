package main

import (
	"sort"
	. "pm/lib/date"
	. "pm/lib/dec2"
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

const catSick CategId_t = 1
const specrefAddon = 2
const segmentStudent = 4

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
	return Select(options).Id(`plancat-`, planId, `-`, categ.categId).Choose(selected)
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

func StateIntOK(state State_t, key string) (int, bool) {
	v := StateValue(state, key)
	if v == `` { return 0, false }
	return Atoi(v), true
}

func StateIntAny(state State_t, keys ...string) (int, bool) {
	for _, key := range keys { if v, ok := StateIntOK(state, key); ok { return v, true } }
	return 0, false
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

func CatPrice(buyYear, age, productId int, categId CategId_t, sickCover int, suppressSurch bool) (base, surcharge EuroCent_t, ok bool) {
	base, surcharge, ok = LookupPrice(buyYear, age, productId)
	if !ok { return 0, 0, false }
	if categId == catSick { base *= EuroCent_t(sickCover / 4500) }
	if suppressSurch { surcharge = 0 }
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

func StateBool(state State_t, keys ...string) bool {
	for _, key := range keys {
		v := Lower(StateValue(state, key))
		switch v {
		case `1`, `on`, `yes`, `true`:
			return true
		}
	}
	return false
}

func CategIs(categ Categ_t, tag string) bool {
	return Contains(Lower(Trim(categ.name)), Lower(Trim(tag)))
}

func PickAnyAddon(choices []CatChoice_t) AddonId_t {
	for _, choice := range choices {
		if choice.addon != 0 { return choice.addon }
	}
	return 0
}

func IsHospDental(categId CategId_t) bool {
	return int(categId) == catHospital || int(categId) == catDental
}

func PlanBundledLevel(plan Plan_t, categId CategId_t) int {
	switch int(categId) {
	case catHospital:
		return plan.hospital
	case catDental:
		return plan.dental
	}
	return 0
}

func LevelLabel(level int) string {
	x, ok := App.lookup.levels.byId[level]
	if !ok { return Str(level) }
	label := Trim(x.label)
	if label != `` { return label }
	return Str(level)
}

func EnsureFreeChoice(plan Plan_t, categId CategId_t, choices []CatChoice_t) []CatChoice_t {
	if !IsHospDental(categId) { return choices }
	for _, choice := range choices {
		if choice.addon == 0 { return choices }
	}

	level := PlanBundledLevel(plan, categId)
	free := CatChoice_t{
		addon: 0,
		level: level,
		label: LevelLabel(level),
	}
	return append([]CatChoice_t{ free }, choices...)
}

func HasZeroChoice(choices []CatChoice_t) bool {
	for _, choice := range choices {
		if choice.addon == 0 { return true }
	}
	return false
}

func HasAddonChoice(choices []CatChoice_t) bool {
	for _, choice := range choices {
		if choice.addon != 0 { return true }
	}
	return false
}

func EnsureOptionalChoice(categ Categ_t, choices []CatChoice_t) []CatChoice_t {
	if categ.required != 0 { return choices }
	if HasZeroChoice(choices) || !HasAddonChoice(choices) { return choices }
	free := CatChoice_t{ addon: 0, level: 0, label: `---` }
	return append([]CatChoice_t{ free }, choices...)
}

func PickRangeChoice(plan Plan_t, categId CategId_t, choices []CatChoice_t, minLevel, maxLevel int) CatChoice_t {
	bundled := PlanBundledLevel(plan, categId)
	if bundled >= minLevel && bundled <= maxLevel {
		return CatChoice_t{ addon: 0, level: bundled, label: LevelLabel(bundled) }
	}
	for _, choice := range choices {
		if choice.addon == 0 { continue }
		if choice.level >= minLevel && choice.level <= maxLevel { return choice }
	}
	for _, choice := range choices {
		if choice.addon != 0 { return choice }
	}
	return CatChoice_t{ addon: 0, level: bundled, label: LevelLabel(bundled) }
}

func CategRangeMatch(plan Plan_t, categId CategId_t, minLevel, maxLevel int) bool {
	bundled := PlanBundledLevel(plan, categId)
	if bundled >= minLevel && bundled <= maxLevel { return true }

	key := PlanCateg_t{ plan: plan.planId, categ: categId }
	for _, choice := range App.lookup.planAddonChoices[key] {
		if choice.addon == 0 { continue }
		if choice.level >= minLevel && choice.level <= maxLevel { return true }
	}
	return false
}

func VisionAmount(plan Plan_t, planBase EuroCent_t) (EuroCent_t, bool) {
	if plan.vision.percent > 0 {
		return ApplyPercent(planBase, plan.vision.percent), true
	}
	if plan.vision.euro > 0 {
		return plan.vision.euro, true
	}
	return 0, false
}

type PlanFilters_t struct {
	segment int
	prior int
	noExam bool
	specref int
	vision bool
	tempVisa bool
	noPVN bool
	naturalMed bool
	deductMin, deductMax int
	hospitalMin, hospitalMax int
	dentalMin, dentalMax int
}

func PlanFilters(state State_t) PlanFilters_t {
	f := PlanFilters_t{
		specref: specrefAddon,
		deductMax: 300000,
		hospitalMax: 999,
		dentalMax: 999,
	}

	if v, ok := StateIntAny(state, `segment`); ok { f.segment = v }
	if v, ok := StateIntAny(state, `priorCov`, `prior`); ok { f.prior = v }
	if v, ok := StateIntAny(state, `noExam`, `exam`); ok { f.noExam = v > 0 }
	if v, ok := StateIntAny(state, `specref`, `referral`); ok { f.specref = v }
	f.vision = StateBool(state, `vision`, `glasses`)
	f.tempVisa = StateBool(state, `tempVisa`, `tempvisa`, `visa`)
	f.noPVN = StateBool(state, `noPVN`, `pvnoff`)
	f.naturalMed = StateBool(state, `naturalMed`, `natural`)

	if v, ok := StateIntAny(state, `deductibleMin`, `deductMin`); ok { f.deductMin = v * 100 }
	if v, ok := StateIntAny(state, `deductibleMax`, `deductMax`); ok { f.deductMax = v * 100 }
	if f.deductMin > f.deductMax { f.deductMin, f.deductMax = f.deductMax, f.deductMin }

	if v, ok := StateIntAny(state, `hospitalMin`, `minHospital`); ok { f.hospitalMin = v }
	if v, ok := StateIntAny(state, `hospitalMax`, `maxHospital`); ok { f.hospitalMax = v }
	if f.hospitalMin > f.hospitalMax { f.hospitalMin, f.hospitalMax = f.hospitalMax, f.hospitalMin }

	if v, ok := StateIntAny(state, `dentalMin`, `minDental`); ok { f.dentalMin = v }
	if v, ok := StateIntAny(state, `dentalMax`, `maxDental`); ok { f.dentalMax = v }
	if f.dentalMin > f.dentalMax { f.dentalMin, f.dentalMax = f.dentalMax, f.dentalMin }

	return f
}

func PlanFilterReasons(plan Plan_t, f PlanFilters_t) []string {
	var reasons []string

	if !f.tempVisa && plan.tempvisa { reasons = append(reasons, `Temp Visa`) }
	if f.segment > 0 && int(plan.segmask)&f.segment == 0 { reasons = append(reasons, `Segment`) }
	if f.prior < plan.priorcov { reasons = append(reasons, `Prior cover`) }
	if f.noExam && f.prior < plan.noexam { reasons = append(reasons, `No exam`) }
	if f.specref != specrefAddon && f.specref != plan.specref && plan.specref != specrefAddon { reasons = append(reasons, `Specialist`) }

	dval := int(plan.ded.adult.euro)
	if dval < f.deductMin || dval > f.deductMax { reasons = append(reasons, `Deductible`) }
	if !CategRangeMatch(plan, CategId_t(catHospital), f.hospitalMin, f.hospitalMax) { reasons = append(reasons, `Hospital`) }
	if !CategRangeMatch(plan, CategId_t(catDental), f.dentalMin, f.dentalMax) { reasons = append(reasons, `Dental`) }

	return reasons
}

type PlanFiltered_t struct {
	planId int
	label string
	reasons []string
}

func FilterRow(x PlanFiltered_t) Elem_t {
	var pills []Elem_t
	for _, reason := range x.reasons {
		pills = append(pills, Span(reason).Class(`filter-pill`))
	}
	return Div().Class(`filter-row`).Id(`fplanId`, x.planId).Wrap(
		Div(x.label).Class(`filter-plan-label`),
		Div().Class(`filter-pill-list`).Wrap(pills),
	)
}

func FilterListCard(filtered []PlanFiltered_t) Elem_t {
	var rows []Elem_t
	for _, x := range filtered { rows = append(rows, FilterRow(x)) }
	if len(rows) == 0 { rows = append(rows, Div(`No plans filtered out.`).Class(`filter-empty`)) }

	return Div().Id(`filterList`).Class(`card`).Wrap(
		Div(`Filtered Out (`, len(filtered), `)`).Class(`card-title`),
		Div().Class(`card-body`, `filter-list-body`).Wrap(
			Elem(`details`).Class(`filter-list-details`).Wrap(
				Elem(`summary`).Class(`filter-list-summary`).Wrap(
					Span(`Show filtered plans`),
				),
				Div().Class(`filter-list-rows`).Wrap(rows),
			),
		),
	)
}

func ListPlans(state State_t) Elem_t {
	buyYear, yearAge, exactAge := PlanAges(state)
	sickCover := StateInt(state, `sickCover`)
	filter := PlanFilters(state)

	var list []Elem_t
	var filtered []PlanFiltered_t
	for planId, plan := range App.lookup.plans.All() {
		age := yearAge
		if plan.exactAge { age = exactAge }
		suppressSurch := !plan.surcharge || filter.segment == segmentStudent
		planBase, planSurch, planOk := CatPrice(buyYear, age, planId, 0, sickCover, suppressSurch)

		reasons := PlanFilterReasons(plan, filter)
		if !planOk || planBase == 0 { reasons = append(reasons, `Zero price`) }
		if len(reasons) > 0 {
			filtered = append(filtered, PlanFiltered_t{
				planId: int(plan.planId),
				label: Str(plan.provName, ` / `, plan.name),
				reasons: reasons,
			})
			continue
		}
		sumBase, sumSurch := EuroCent_t(0), EuroCent_t(0)
		if planOk { sumBase += planBase; sumSurch += planSurch }

		addonRows := []Elem_t{
			PlanAddonHead(),
			PlanAddonRow(`Plan`, ``, PriceText(planBase, planOk), PriceText(planSurch, planOk)),
		}

		planKey := plan.planId
			for _, categ := range App.lookup.categs.All() {
				if categ.categId == 0 { continue }
				if filter.noPVN && CategIs(categ, `pvn`) { continue }

				key := PlanCateg_t{ plan: planKey, categ: categ.categId }
				choices := App.lookup.planAddonChoices[key]
				if len(choices) == 0 { continue }
				choices = EnsureFreeChoice(plan, categ.categId, choices)
			choices = EnsureOptionalChoice(categ, choices)

			choice, ok := App.lookup.planAddons[key]
			if !ok { continue }
				if IsHospDental(categ.categId) {
					minLevel, maxLevel := filter.hospitalMin, filter.hospitalMax
					if int(categ.categId) == catDental { minLevel, maxLevel = filter.dentalMin, filter.dentalMax }
					choice = PickRangeChoice(plan, categ.categId, choices, minLevel, maxLevel)
				}
				if CategIs(categ, `natural`) {
					if !filter.naturalMed {
						choice.addon = 0
						choice.label = `---`
					} else if choice.addon == 0 {
						choice.addon = PickAnyAddon(choices)
					}
				}
			hasMulti := len(choices) > 1
			if !hasMulti && choice.addon == 0 { continue }

			var base, surch EuroCent_t
			priceOk := false
			if choice.addon != 0 {
					base, surch, priceOk = CatPrice(buyYear, age, int(choice.addon), categ.categId, sickCover, suppressSurch)
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
		if filter.vision {
			visionBase, visionOk := VisionAmount(plan, planBase)
			if visionOk {
				sumBase += visionBase
				addonRows = append(addonRows,
					PlanAddonRow(`Vision`, ``, PriceText(visionBase, true), `-`),
				)
			}
		}

		sumBaseText := PriceText(sumBase, true)
		sumSurchText := PriceText(sumSurch, true)
		totalText := PriceText(sumBase + sumSurch, true)

		list = append(list, PlanCard(planId, plan, totalText, sumBaseText, sumSurchText, addonRows))
	}

	sort.Slice(filtered, func(i, j int) bool {
		return Lower(filtered[i].label) < Lower(filtered[j].label)
	})

	planList := Div().Id(`planList`).Class(`card`).Wrap(
		Div(`Plans (`, len(list), `)`).Class(`card-title`),
		Div().Class(`card-body`, `plan-list-body`).Wrap(list),
	)
	filterList := FilterListCard(filtered)

	return Div().Id(`PlansAndFilters`).Wrap(
		planList,
		filterList,
	)
}
