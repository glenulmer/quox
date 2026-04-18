package main

import (
	"fmt"
	"sort"
	"strings"

	. "quo2/lib/dec2"
	. "quo2/lib/output"
)

const editQDepMaxCount = 5
const editQPrimeModePct = `pct`
const editQPrimeModeEur = `eur`

type EditQDep_t struct {
	depId int
	name string
	birth string
	vision bool
}

type EditQPrimeCharge_t struct {
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

func EditQPrimeModeKey(itemId int, categId CategId_t) string {
	return Str(`editq-prime-`, itemId, `-`, categId, `-mode`)
}

func EditQPrimeAmountKey(itemId int, categId CategId_t) string {
	return Str(`editq-prime-`, itemId, `-`, categId, `-amount`)
}

func EditQPrimeNoteKey(itemId int, categId CategId_t) string {
	return Str(`editq-prime-`, itemId, `-`, categId, `-note`)
}

func EditQPrimeModeControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-prime-%d-%d-mode`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPrimeAmountControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-prime-%d-%d-amount`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPrimeNoteControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-prime-%d-%d-note`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQDepChargeModeKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-prex-`, itemId, `-`, categId, `-mode`)
}

func EditQDepChargeAmountKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-prex-`, itemId, `-`, categId, `-amount`)
}

func EditQDepChargeNoteKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-prex-`, itemId, `-`, categId, `-note`)
}

func EditQDepChargeModeControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-prex-%d-%d-mode`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeAmountControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-prex-%d-%d-amount`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeNoteControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-prex-%d-%d-note`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQPrimeMode(v string) string {
	if Lower(Trim(v)) == editQPrimeModeEur { return editQPrimeModeEur }
	return editQPrimeModePct
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

func EditQPrimeAppliedAmount(mode, amount string, base EuroCent_t) EuroCent_t {
	n, ok := EditQParseDecimal100(amount)
	if !ok || n <= 0 { return 0 }
	if EditQPrimeMode(mode) == editQPrimeModeEur {
		return EuroCent_t(n)
	}
	// n is percent*100, apply to base euro-cent amount.
	return EuroCent_t((int64(base) * n) / 10000)
}

func EditQPrimeCharges(vars UIBagVars_t) []EditQPrimeCharge_t {
	state := QuoteStateFromVars(vars)
	selected := QuoteSelectedRows(state)
	var out []EditQPrimeCharge_t
	for _, x := range selected {
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			modeKey := EditQPrimeModeKey(x.item.itemId, 0)
			mode := EditQPrimeMode(vars[modeKey])
			amount := vars[EditQPrimeAmountKey(x.item.itemId, 0)]
			applied := EditQPrimeAppliedAmount(mode, amount, x.row.planBase)
			out = append(out, EditQPrimeCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: 0,
				level: `Plan`,
				base: x.row.planBase,
				mode: mode,
				amount: amount,
				note: vars[EditQPrimeNoteKey(x.item.itemId, 0)],
				applied: applied,
			})
		}

		for _, addon := range x.row.addons {
			if !addon.priceOk { continue }
			if addon.base+addon.surcharge <= 0 { continue }
			modeKey := EditQPrimeModeKey(x.item.itemId, addon.categId)
			mode := EditQPrimeMode(vars[modeKey])
			amount := vars[EditQPrimeAmountKey(x.item.itemId, addon.categId)]
			applied := EditQPrimeAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPrimeCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: addon.categId,
				level: level,
				base: addon.base,
				mode: mode,
				amount: amount,
				note: vars[EditQPrimeNoteKey(x.item.itemId, addon.categId)],
				applied: applied,
			})
		}
	}
	return out
}

func EditQDependentCharges(vars UIBagVars_t, dep EditQDep_t) []EditQPrimeCharge_t {
	work := QuoteStateFromVars(vars)
	work.quote[`birth`] = dep.birth
	if dep.vision {
		work.quote[`vision`] = `1`
	} else {
		work.quote[`vision`] = ``
	}
	_, yearAge, exactAge := PlanAges(work)
	selected := QuoteSelectedRows(work)

	var out []EditQPrimeCharge_t
	for _, x := range selected {
		age := yearAge
		ageMode := `year`
		if EditQPlanUsesExactAge(x.row.planId) {
			age = exactAge
			ageMode = `exact`
		}
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			modeKey := EditQDepChargeModeKey(dep.depId, x.item.itemId, 0)
			mode := EditQPrimeMode(vars[modeKey])
			amount := vars[EditQDepChargeAmountKey(dep.depId, x.item.itemId, 0)]
			applied := EditQPrimeAppliedAmount(mode, amount, x.row.planBase)
			out = append(out, EditQPrimeCharge_t{
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
			if addon.base+addon.surcharge <= 0 { continue }
			modeKey := EditQDepChargeModeKey(dep.depId, x.item.itemId, addon.categId)
			mode := EditQPrimeMode(vars[modeKey])
			amount := vars[EditQDepChargeAmountKey(dep.depId, x.item.itemId, addon.categId)]
			applied := EditQPrimeAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPrimeCharge_t{
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

func EditQCurrentYear() int {
	year := 0
	row := App.DB.CallRow(`klec_current_year_query`).Scan(&year)
	if !row.HasError() && year > 0 { return year }
	today := CurrentDBDate()
	if today.Year() > 0 { return today.Year() }
	return 2026
}

func EditQDefaultDependentBirth() string {
	return fmt.Sprintf(`%04d-06-15`, EditQCurrentYear()-4)
}

func EditQBirthSortKey(v string) string {
	v = Trim(v)
	if len(v) != len(`2006-01-02`) { return `9999-99-99` }
	if v[4] != '-' || v[7] != '-' { return `9999-99-99` }
	return v
}

func EditQDependents(vars UIBagVars_t, sortForGet bool) []EditQDep_t {
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
		if x.birth == `` { x.birth = EditQDefaultDependentBirth() }
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

func EditQDeleteDependent(vars UIBagVars_t, depId int) {
	prefix := Str(`editq-dep-`, depId, `-`)
	for key := range vars {
		if strings.HasPrefix(key, prefix) { delete(vars, key) }
	}
}

func EditQAddDependent(state *State_t) bool {
	deps := EditQDependents(state.quote, false)
	if len(deps) >= editQDepMaxCount { return false }
	next := 1
	for _, dep := range deps {
		if dep.depId >= next { next = dep.depId + 1 }
	}
	state.quote[EditQDepNameKey(next)] = ``
	state.quote[EditQDepBirthKey(next)] = EditQDefaultDependentBirth()
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

	if itemId, categId, ok := EditQPrimeModeControl(name); ok {
		state.quote[EditQPrimeModeKey(itemId, categId)] = EditQPrimeMode(value)
		return true
	}
	if _, _, ok := EditQPrimeAmountControl(name); ok {
		state.quote[name] = value
		return true
	}
	if _, _, ok := EditQPrimeNoteControl(name); ok {
		state.quote[name] = value
		return true
	}
	if depId, itemId, categId, ok := EditQDepChargeModeControl(name); ok {
		state.quote[EditQDepChargeModeKey(depId, itemId, categId)] = EditQPrimeMode(value)
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
		EditQAddDependent(state)
		return true
	}
	if depId, ok := EditQDepDelControl(name); ok {
		EditQDeleteDependent(state.quote, depId)
		if len(EditQDependents(state.quote, false)) == 0 {
			EditQAddDependent(state)
		}
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

func EditQEnsureDefaultDependent(state *State_t) {
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	if len(EditQDependents(state.quote, false)) > 0 { return }
	EditQAddDependent(state)
}
