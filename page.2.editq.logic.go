package main

import (
	"fmt"
	"sort"
	"strings"

	. "pm/lib/dec2"
	. "pm/lib/output"
)

const editQDepMaxCount = 5
const editQPreSeqKey = `editq-pre-seq`
const editQDepSeqKey = `editq-dep-seq`
const editQPrimeModePct = `pct`
const editQPrimeModeEur = `eur`

type EditQCond_t struct {
	condId int
	text string
}

type EditQDep_t struct {
	depId int
	name string
	birth string
	vision bool
	conds []EditQCond_t
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

func EditQPreKey(condId int) string { return Str(`editq-pre-`, condId) }

func EditQPreControl(name string) (condId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-pre-%d`, &condId)
	if err != nil || n != 1 || condId <= 0 { return 0, false }
	return condId, true
}

func EditQPreAddControlName() string { return `editq-pre-add` }
func EditQPreAddControl(name string) bool { return name == EditQPreAddControlName() }

func EditQPreDelControlName(condId int) string { return Str(`editq-pre-del-`, condId) }
func EditQPreDelControl(name string) (condId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-pre-del-%d`, &condId)
	if err != nil || n != 1 || condId <= 0 { return 0, false }
	return condId, true
}

func EditQDepNameKey(depId int) string { return Str(`editq-dep-`, depId, `-name`) }
func EditQDepBirthKey(depId int) string { return Str(`editq-dep-`, depId, `-birth`) }
func EditQDepVisionKey(depId int) string { return Str(`editq-dep-`, depId, `-vision`) }
func EditQDepCondSeqKey(depId int) string { return Str(`editq-dep-`, depId, `-cond-seq`) }
func EditQDepCondKey(depId, condId int) string { return Str(`editq-dep-`, depId, `-cond-`, condId) }

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

func EditQDepCondControl(name string) (depId, condId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-cond-%d`, &depId, &condId)
	if err != nil || n != 2 || depId <= 0 || condId <= 0 { return 0, 0, false }
	return depId, condId, true
}

func EditQDepAddControlName() string { return `editq-dep-add` }
func EditQDepAddControl(name string) bool { return name == EditQDepAddControlName() }

func EditQDepDelControlName(depId int) string { return Str(`editq-dep-del-`, depId) }
func EditQDepDelControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-del-%d`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepPreAddControlName(depId int) string { return Str(`editq-dep-pre-add-`, depId) }
func EditQDepPreAddControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-pre-add-%d`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepPreDelControlName(depId, condId int) string { return Str(`editq-dep-pre-del-`, depId, `-`, condId) }
func EditQDepPreDelControl(name string) (depId, condId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-pre-del-%d-%d`, &depId, &condId)
	if err != nil || n != 2 || depId <= 0 || condId <= 0 { return 0, 0, false }
	return depId, condId, true
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

func EditQPrimeCharges(vars QuoteVars_t) []EditQPrimeCharge_t {
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

func EditQDependentCharges(vars QuoteVars_t, dep EditQDep_t) []EditQPrimeCharge_t {
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

func EditQDependentAgeText(vars QuoteVars_t, dep EditQDep_t) string {
	work := QuoteStateFromVars(vars)
	work.quote[`birth`] = dep.birth
	buyYear, yearAge, exactAge := PlanAges(work)
	_ = buyYear

	selected := QuoteSelectedRows(work)
	if len(selected) == 0 {
		if yearAge > 0 { return Str(`age `, yearAge) }
		return ``
	}

	hasExact := false
	hasYear := false
	for _, x := range selected {
		if EditQPlanUsesExactAge(x.row.planId) {
			hasExact = true
			continue
		}
		hasYear = true
	}

	if hasExact && hasYear {
		if exactAge > 0 && yearAge > 0 { return Str(`age `, exactAge, `/`, yearAge) }
	}
	if hasExact {
		if exactAge > 0 { return Str(`age `, exactAge) }
		return ``
	}
	if yearAge > 0 { return Str(`age `, yearAge) }
	return ``
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

func EditQPreConditions(vars QuoteVars_t) []EditQCond_t {
	var out []EditQCond_t
	for key, value := range vars {
		condId, ok := EditQPreControl(key)
		if !ok { continue }
		out = append(out, EditQCond_t{ condId:condId, text:value })
	}
	sort.Slice(out, func(i, j int) bool { return out[i].condId < out[j].condId })
	return out
}

func EditQBirthSortKey(v string) string {
	v = Trim(v)
	if len(v) != len(`2006-01-02`) { return `9999-99-99` }
	if v[4] != '-' || v[7] != '-' { return `9999-99-99` }
	return v
}

func EditQDependents(vars QuoteVars_t, sortForGet bool) []EditQDep_t {
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
		depId, condId, ok := EditQDepCondControl(key)
		if !ok { continue }
		x := all[depId]
		x.depId = depId
		x.conds = append(x.conds, EditQCond_t{ condId:condId, text:value })
		all[depId] = x
	}

	var out []EditQDep_t
	for _, x := range all {
		if x.birth == `` { x.birth = EditQDefaultDependentBirth() }
		sort.Slice(x.conds, func(i, j int) bool { return x.conds[i].condId < x.conds[j].condId })
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

func EditQDeleteDependent(vars QuoteVars_t, depId int) {
	prefix := Str(`editq-dep-`, depId, `-`)
	for key := range vars {
		if strings.HasPrefix(key, prefix) { delete(vars, key) }
	}
}

func EditQAddDependent(state *State_t) bool {
	deps := EditQDependents(state.quote, false)
	if len(deps) >= editQDepMaxCount { return false }
	next := StateInt(*state, editQDepSeqKey) + 1
	state.quote[editQDepSeqKey] = Str(next)
	state.quote[EditQDepNameKey(next)] = ``
	state.quote[EditQDepBirthKey(next)] = EditQDefaultDependentBirth()
	state.quote[EditQDepVisionKey(next)] = ``
	state.quote[EditQDepCondSeqKey(next)] = `0`
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

	if EditQPreAddControl(name) {
		next := StateInt(*state, editQPreSeqKey) + 1
		state.quote[editQPreSeqKey] = Str(next)
		state.quote[EditQPreKey(next)] = ``
		return true
	}
	if condId, ok := EditQPreDelControl(name); ok {
		delete(state.quote, EditQPreKey(condId))
		return true
	}
	if _, ok := EditQPreControl(name); ok {
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

	if depId, ok := EditQDepPreAddControl(name); ok {
		seqKey := EditQDepCondSeqKey(depId)
		next := Atoi(state.quote[seqKey]) + 1
		state.quote[seqKey] = Str(next)
		state.quote[EditQDepCondKey(depId, next)] = ``
		return true
	}
	if depId, condId, ok := EditQDepPreDelControl(name); ok {
		delete(state.quote, EditQDepCondKey(depId, condId))
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
	if _, _, ok := EditQDepCondControl(name); ok {
		state.quote[name] = value
		return true
	}

	return false
}

func EditQEnsureFirstPlanSelected(state *State_t) {
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	if len(QuoteSelectedItems(state.quote)) > 0 { return }
	plans := QuotePlans(*state).plans
	if len(plans) == 0 { return }
	QuoteSelectedAdd(state, int(plans[0].planId))
}

func EditQEnsureDefaultDependent(state *State_t) {
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	if len(EditQDependents(state.quote, false)) > 0 { return }
	EditQAddDependent(state)
}
