package main

import (
	"fmt"
	"sort"
	"strings"

	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

const editQDepMaxCount = 5
const editQPreexModePct = `pct`
const editQPreexModeEur = `eur`

type EditQDep_t struct {
	depId int
	name string
	birth string
	vision bool
}

type EditQPreexCharge_t struct {
	itemId int
	plan string
	planPrice EuroCent_t
	planAge int
	planAgeMode string
	categId CategId_t
	level string
	base EuroCent_t
	mode string
	amount string
	note string
	applied EuroCent_t
}

func EditQDepNameKey(depId int) string { return Str(`editq-dep-`, depId, `-name`) }
func EditQDepBirthKey(depId int) string { return Str(`editq-dep-`, depId, `-birth`) }
func EditQDepVisionKey(depId int) string { return Str(`editq-dep-`, depId, `-vision`) }

func EditQDepNameControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-name`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepBirthControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-birth`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepVisionControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-vision`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepAddControlName() string { return `editq-dep-add` }
func EditQDepAddControl(name string) bool { return name == EditQDepAddControlName() }

func EditQDepDelControlName(depId int) string { return Str(`editq-dep-del-`, depId) }
func EditQDepDelControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-del-%d`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQPreexModeKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-mode`)
}

func EditQPreexAmountKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-amount`)
}

func EditQPreexNoteKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-note`)
}

func EditQPreexModeControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-mode`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPreexAmountControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-amount`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPreexNoteControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-note`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQDepChargeModeKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-mode`)
}

func EditQDepChargeAmountKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-amount`)
}

func EditQDepChargeNoteKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-note`)
}

func EditQDepChargeModeControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-mode`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeAmountControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-amount`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeNoteControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-note`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQPreexMode(v string) string {
	if Lower(Trim(v)) == editQPreexModeEur { return editQPreexModeEur }
	return editQPreexModePct
}

func EditQParseDecimal100(v string) (int64, bool) {
	s := Trim(v)
	s = strings.ReplaceAll(s, `%`, ``)
	s = strings.ReplaceAll(s, `€`, ``)
	s = strings.ReplaceAll(s, `eur`, ``)
	s = strings.ReplaceAll(s, ` `, ``)
	s = strings.ReplaceAll(s, `.`, `,`)
	if s == `` { return 0, false }
	parts := Split(s, `,`)
	if len(parts) > 2 { return 0, false }
	onlyDigits := func(in string) string {
		var out []byte
		for i := 0; i < len(in); i++ {
			if in[i] >= '0' && in[i] <= '9' { out = append(out, in[i]) }
		}
		return string(out)
	}
	whole := onlyDigits(parts[0])
	if whole == `` { whole = `0` }
	frac := ``
	if len(parts) == 2 { frac = onlyDigits(parts[1]) }
	if len(frac) > 2 { frac = frac[:2] }
	for len(frac) < 2 { frac += `0` }
	return int64(Atoi(whole)*100 + Atoi(frac)), true
}

func EditQPreexAppliedAmount(mode, amount string, base EuroCent_t) EuroCent_t {
	n, ok := EditQParseDecimal100(amount)
	if !ok || n <= 0 { return 0 }
	if EditQPreexMode(mode) == editQPreexModeEur {
		return EuroCent_t(n)
	}
	// n is percent*100, apply to base euro-cent amount.
	return EuroCent_t((int64(base) * n) / 10000)
}

func EditQPreexCharges(vars UIBagVars_t) []EditQPreexCharge_t {
	state := QuoteStateFromVars(vars)
	selected := QuoteSelectedRows(state)
	var out []EditQPreexCharge_t
	for _, x := range selected {
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			modeKey := EditQPreexModeKey(x.item.itemId, 0)
			mode := EditQPreexMode(vars[modeKey])
			amount := vars[EditQPreexAmountKey(x.item.itemId, 0)]
			applied := EditQPreexAppliedAmount(mode, amount, x.row.planBase)
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: 0,
				level: `Plan`,
				base: x.row.planBase,
				mode: mode,
				amount: amount,
				note: vars[EditQPreexNoteKey(x.item.itemId, 0)],
				applied: applied,
			})
		}

		for _, addon := range x.row.addons {
			if !addon.priceOk { continue }
			if addon.categId <= 0 { continue }
			if Contains(Lower(Trim(addon.categ)), `vision`) { continue }
			if addon.base+addon.surcharge <= 0 { continue }
			modeKey := EditQPreexModeKey(x.item.itemId, addon.categId)
			mode := EditQPreexMode(vars[modeKey])
			amount := vars[EditQPreexAmountKey(x.item.itemId, addon.categId)]
			applied := EditQPreexAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: addon.categId,
				level: level,
				base: addon.base,
				mode: mode,
				amount: amount,
				note: vars[EditQPreexNoteKey(x.item.itemId, addon.categId)],
				applied: applied,
			})
		}
	}
	return out
}

func EditQDependantCharges(vars UIBagVars_t, dep EditQDep_t) []EditQPreexCharge_t {
	work := QuoteStateFromVars(vars)
	work.quote[`birth`] = dep.birth
	if dep.vision {
		work.quote[`vision`] = `1`
	} else {
		work.quote[`vision`] = ``
	}
	_, yearAge, exactAge := PlanAges(work)
	selected := QuoteSelectedRows(work)

	var out []EditQPreexCharge_t
	for _, x := range selected {
		age := yearAge
		ageMode := `year`
		if EditQPlanUsesExactAge(x.row.planId) {
			age = exactAge
			ageMode = `exact`
		}
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			modeKey := EditQDepChargeModeKey(dep.depId, x.item.itemId, 0)
			mode := EditQPreexMode(vars[modeKey])
			amount := vars[EditQDepChargeAmountKey(dep.depId, x.item.itemId, 0)]
			applied := EditQPreexAppliedAmount(mode, amount, x.row.planBase)
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				planAge: age,
				planAgeMode: ageMode,
				categId: 0,
				level: `Plan`,
				base: x.row.planBase,
				mode: mode,
				amount: amount,
				note: vars[EditQDepChargeNoteKey(dep.depId, x.item.itemId, 0)],
				applied: applied,
			})
		}
		for _, addon := range x.row.addons {
			if !addon.priceOk { continue }
			if addon.categId <= 0 { continue }
			if Contains(Lower(Trim(addon.categ)), `vision`) { continue }
			if addon.base+addon.surcharge <= 0 { continue }
			modeKey := EditQDepChargeModeKey(dep.depId, x.item.itemId, addon.categId)
			mode := EditQPreexMode(vars[modeKey])
			amount := vars[EditQDepChargeAmountKey(dep.depId, x.item.itemId, addon.categId)]
			applied := EditQPreexAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				planAge: age,
				planAgeMode: ageMode,
				categId: addon.categId,
				level: level,
				base: addon.base,
				mode: mode,
				amount: amount,
				note: vars[EditQDepChargeNoteKey(dep.depId, x.item.itemId, addon.categId)],
				applied: applied,
			})
		}
	}
	return out
}

func EditQPlanUsesExactAge(planId int) bool {
	for id, plan := range App.lookup.plans.All() {
		if id != planId { continue }
		return plan.exactAge
	}
	return false
}

func EditQDefaultDependantBirth(vars UIBagVars_t) string {
	year := 0
	buy := QuoteParseBuyDate(vars[`buy`])
	if buy.Year() > 0 { year = buy.Year() }
	if year <= 0 {
		_, _, def := QuoteBuyBounds()
		if def.Year() > 0 { year = def.Year() }
	}
	if year <= 0 {
		today := CurrentDBDate()
		if today.Year() > 0 { year = today.Year() }
	}
	if year <= 0 { year = 2026 }
	return fmt.Sprintf(`%04d-06-15`, year-4)
}

func EditQBirthSortKey(v string) string {
	v = Trim(v)
	if len(v) != len(`2006-01-02`) { return `9999-99-99` }
	if v[4] != '-' || v[7] != '-' { return `9999-99-99` }
	return v
}

func EditQDependants(vars UIBagVars_t, sortForGet bool) []EditQDep_t {
	all := make(map[int]EditQDep_t)
	for key, value := range vars {
		if depId, ok := EditQDepNameControl(key); ok {
			x := all[depId]
			x.depId = depId
			x.name = value
			all[depId] = x
			continue
		}
		if depId, ok := EditQDepBirthControl(key); ok {
			x := all[depId]
			x.depId = depId
			x.birth = value
			all[depId] = x
			continue
		}
		if depId, ok := EditQDepVisionControl(key); ok {
			x := all[depId]
			x.depId = depId
			x.vision = QuoteVarBool(value)
			all[depId] = x
			continue
		}
	}

	var out []EditQDep_t
	for _, x := range all {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].depId != out[j].depId { return out[i].depId < out[j].depId }
		return false
	})
	if !sortForGet { return out }

	sort.Slice(out, func(i, j int) bool {
		bi, bj := EditQBirthSortKey(out[i].birth), EditQBirthSortKey(out[j].birth)
		if bi != bj { return bi < bj }
		ni, nj := Lower(Trim(out[i].name)), Lower(Trim(out[j].name))
		if ni != nj { return ni < nj }
		return out[i].depId < out[j].depId
	})
	return out
}

func EditQDeleteDependant(vars UIBagVars_t, depId int) {
	prefix := Str(`editq-dep-`, depId, `-`)
	for key := range vars {
		if strings.HasPrefix(key, prefix) { delete(vars, key) }
	}
}

func EditQAddDependant(state *State_t) bool {
	deps := EditQDependants(state.quote, false)
	if len(deps) >= editQDepMaxCount { return false }
	next := 1
	for _, dep := range deps {
		if dep.depId >= next { next = dep.depId + 1 }
	}
	state.quote[EditQDepNameKey(next)] = ``
	state.quote[EditQDepBirthKey(next)] = EditQDefaultDependantBirth(state.quote)
	state.quote[EditQDepVisionKey(next)] = ``
	return true
}

func EditQApply(state *State_t, name, value string) bool {
	if name == `` { return false }
	if state.quote == nil { state.quote = QuoteDefaultVars() }

	if name == `clientName` {
		state.quote[name] = value
		return true
	}
	if name == `lang` {
		if Atoi(value) <= 0 { value = Str(int(English)) }
		state.quote[name] = value
		return true
	}
	if name == `slim` {
		if Atoi(value) == 1 { state.quote[name] = `1` } else { state.quote[name] = `0` }
		return true
	}

	if itemId, categId, ok := EditQPreexModeControl(name); ok {
		state.quote[EditQPreexModeKey(itemId, categId)] = EditQPreexMode(value)
		return true
	}
	if _, _, ok := EditQPreexAmountControl(name); ok {
		state.quote[name] = value
		return true
	}
	if _, _, ok := EditQPreexNoteControl(name); ok {
		state.quote[name] = value
		return true
	}
	if depId, itemId, categId, ok := EditQDepChargeModeControl(name); ok {
		state.quote[EditQDepChargeModeKey(depId, itemId, categId)] = EditQPreexMode(value)
		return true
	}
	if _, _, _, ok := EditQDepChargeAmountControl(name); ok {
		state.quote[name] = value
		return true
	}
	if _, _, _, ok := EditQDepChargeNoteControl(name); ok {
		state.quote[name] = value
		return true
	}

	if EditQDepAddControl(name) {
		EditQAddDependant(state)
		return true
	}
	if depId, ok := EditQDepDelControl(name); ok {
		EditQDeleteDependant(state.quote, depId)
		return true
	}

	if _, ok := EditQDepNameControl(name); ok {
		state.quote[name] = value
		return true
	}
	if _, ok := EditQDepBirthControl(name); ok {
		state.quote[name] = value
		return true
	}
	if _, ok := EditQDepVisionControl(name); ok {
		state.quote[name] = value
		return true
	}
	return false
}

func EditQDropPreinsertedDependant(state *State_t) {
	if state == nil { return }
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	deps := EditQDependants(state.quote, false)
	if len(deps) != 1 { return }

	dep := deps[0]
	if dep.depId <= 0 { return }
	if Trim(dep.name) != `` || dep.vision { return }
	if Trim(dep.birth) != EditQDefaultDependantBirth(state.quote) { return }
	prefix := Str(`editq-dep-`, dep.depId, `-preex-`)
	for key, value := range state.quote {
		if !strings.HasPrefix(key, prefix) { continue }
		if Trim(value) != `` { return }
	}
	EditQDeleteDependant(state.quote, dep.depId)
}
