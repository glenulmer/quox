package main

import (
	"strconv"
	"strings"
)

const forcedCatNatural CategId_t = 6

var forcedQuoteDefaults bool

func SetForcedQuoteDefaults() { forcedQuoteDefaults = true }

func ForcedPlan(provider, name string) int {
	provider = strings.TrimSpace(strings.ToLower(provider))
	name = strings.TrimSpace(strings.ToLower(name))
	for id, plan := range App.lookup.plans.All() {
		if strings.TrimSpace(strings.ToLower(plan.provName)) != provider { continue }
		if strings.TrimSpace(strings.ToLower(plan.name)) != name { continue }
		return id
	}
	return 0
}

func ForcedSHIPlan() int {
	if x := ForcedPlan(`SHI`, `KV`); x > 0 { return x }
	best := 0
	for id, plan := range App.lookup.plans.All() {
		if strings.TrimSpace(strings.ToLower(plan.provName)) != strings.ToLower(`SHI`) { continue }
		if best == 0 || id < best { best = id }
	}
	return best
}

func ForcedAddon(planId int, categId CategId_t, labels ...string) (AddonId_t, bool) {
	key := PlanCateg_t{ plan: PlanId_t(planId), categ: categId }
	choices := App.lookup.planAddonChoices[key]
	if len(choices) == 0 { return 0, false }
	for _, want := range labels {
		want = strings.TrimSpace(strings.ToLower(want))
		for _, choice := range choices {
			if strings.TrimSpace(strings.ToLower(choice.label)) == want {
				return choice.addon, true
			}
		}
	}
	return 0, false
}

func ForcedAdd(state *State_t, planId int) int {
	if state == nil || state.quote == nil || planId <= 0 { return 0 }
	before, _ := strconv.Atoi(state.quote[quoteSelectedSeqKey])
	QuoteSelectedAdd(state, planId)
	after, _ := strconv.Atoi(state.quote[quoteSelectedSeqKey])
	if after <= before { return 0 }
	return after
}

func ForcedSet(vars UIBagVars_t, itemId int, categId CategId_t, addon AddonId_t) {
	if vars == nil || itemId <= 0 { return }
	vars[QuoteSelectedCatKey(itemId, categId)] = strconv.Itoa(int(addon))
}

func ForcedPick(vars UIBagVars_t, itemId int, planId int, categId CategId_t, labels ...string) bool {
	addon, ok := ForcedAddon(planId, categId, labels...)
	if !ok { return false }
	ForcedSet(vars, itemId, categId, addon)
	return true
}

func ForcedPickLevel(vars UIBagVars_t, itemId int, planId int, categId CategId_t, level int) bool {
	plan, ok := App.lookup.plans.byId[planId]
	if ok {
		if categId == CategId_t(catHospital) && plan.hospital == level {
			ForcedSet(vars, itemId, categId, 0)
			return true
		}
		if categId == CategId_t(catDental) && plan.dental == level {
			ForcedSet(vars, itemId, categId, 0)
			return true
		}
	}

	key := PlanCateg_t{ plan: PlanId_t(planId), categ: categId }
	for _, choice := range App.lookup.planAddonChoices[key] {
		if choice.level != level { continue }
		ForcedSet(vars, itemId, categId, choice.addon)
		return true
	}
	return false
}

func QuoteApplyForcedQuoteDefaults(vars UIBagVars_t) {
	if !forcedQuoteDefaults || vars == nil { return }

	vars[`clientName`] = `Jill Jones`
	vars[`birth`] = `1994-01-15`
	vars[`vision`] = `1`
	vars[`sickCover`] = `80000`

	QuoteDropKeysByPrefix(vars, `selplan-`)
	QuoteDropKeysByPrefix(vars, `selcat-`)
	vars[quoteSelectedSeqKey] = `0`

	vars[`naturalMed`] = `1`

	state := State_t{ quote: vars }
	inter := ForcedPlan(`Inter`, `LA-VNS U`)
	gothaer := ForcedPlan(`Gothaer`, `MediVita250`)
	bbkk := ForcedPlan(`BBKK`, `GesundheitVARIO 800`)
	shi := ForcedSHIPlan()

	if inter > 0 {
		item := ForcedAdd(&state, inter)
		_ = ForcedPick(vars, item, inter, catSick, `43A`)
	}
	if gothaer > 0 {
		item := ForcedAdd(&state, gothaer)
		_ = ForcedPick(vars, item, gothaer, catSick, `43A`)
		_ = ForcedPickLevel(vars, item, gothaer, CategId_t(catHospital), 31) // Ward
		ForcedSet(vars, item, CategId_t(catDental), 0) // Cancel / No Dental
		_ = ForcedPick(vars, item, gothaer, forcedCatNatural, `Natural`)
	}
	if bbkk > 0 {
		item := ForcedAdd(&state, bbkk)
		_ = ForcedPickLevel(vars, item, bbkk, CategId_t(catHospital), 31) // Ward
		_ = ForcedPickLevel(vars, item, bbkk, CategId_t(catDental), 42) // Better
	}
	if shi > 0 {
		_ = ForcedAdd(&state, shi)
	}

	QuoteDropKeysByPrefix(vars, `editq-preex-`)
	QuoteDropKeysByPrefix(vars, `editq-dep-`)

	selected := QuoteSelectedItems(vars)
	for _, item := range selected {
		switch item.planId {
		case inter:
			vars[EditQPreexModeKey(item.itemId, 0)] = editQPreexModeEur
			vars[EditQPreexAmountKey(item.itemId, 0)] = `9,34`
		case gothaer:
			vars[EditQPreexModeKey(item.itemId, 0)] = editQPreexModeEur
			vars[EditQPreexAmountKey(item.itemId, 0)] = `35,47`
		case bbkk:
			vars[EditQPreexModeKey(item.itemId, 0)] = editQPreexModeEur
			vars[EditQPreexAmountKey(item.itemId, 0)] = `50`
		case shi:
			vars[EditQPreexModeKey(item.itemId, 2)] = editQPreexModeEur
			vars[EditQPreexAmountKey(item.itemId, 2)] = `55`
		}
	}

	vars[EditQDepNameKey(1)] = `Melvin`
	vars[EditQDepBirthKey(1)] = `2018-07-01`
	vars[EditQDepVisionKey(1)] = ``
}

func QuoteJaneBirthDate() string {
	birth := EditQDefaultDependentBirth()
	if len(birth) != len(`2006-01-02`) { return birth }
	return birth[:5] + `07-15`
}

func QuoteDropKeysByPrefix(vars UIBagVars_t, prefix string) {
	for key := range vars {
		if strings.HasPrefix(key, prefix) { delete(vars, key) }
	}
}
