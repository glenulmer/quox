package main

import (
	"fmt"
	"sort"

	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

const quoteSelectedMaxCount = 5

type QuoteSelectedItem_t struct {
	itemId int
	planId int
	cats map[CategId_t]AddonId_t
}

type QuoteSelectedRow_t struct {
	item QuoteSelectedItem_t
	row QuotePlan_t
}

func QuoteSelectedTitle(count int) string {
	if count < 0 { count = 0 }
	if count > quoteSelectedMaxCount { count = quoteSelectedMaxCount }
	return Str(`Selected Plans (` , count, ` / `, quoteSelectedMaxCount, `)`)
}

func QuoteSelectedPlanControl(name string) (itemId int, ok bool) {
	n, err := fmt.Sscanf(name, `selplan-%d`, &itemId)
	if err != nil || n != 1 || itemId <= 0 { return 0, false }
	return itemId, true
}

func QuoteSelectedCatKey(itemId int, categId CategId_t) string {
	return Str(`selcat-`, itemId, `-`, categId)
}

func QuoteSelectedCatControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `selcat-%d-%d`, &itemId, &cat)
	if err != nil || n != 2 { return 0, 0, false }
	if itemId <= 0 || cat <= 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func QuoteSelectedAddControlName(planId int) string {
	return Str(`seladd-`, planId)
}

func QuoteSelectedAddControl(name string) (planId int, ok bool) {
	n, err := fmt.Sscanf(name, `seladd-%d`, &planId)
	if err != nil || n != 1 || planId <= 0 { return 0, false }
	return planId, true
}

func QuoteSelectedDelControlName(itemId int) string {
	return Str(`seldel-`, itemId)
}

func QuoteSelectedDelControl(name string) (itemId int, ok bool) {
	n, err := fmt.Sscanf(name, `seldel-%d`, &itemId)
	if err != nil || n != 1 || itemId <= 0 { return 0, false }
	return itemId, true
}

func QuoteSelectedItems(vars QuoteVars_t) []QuoteSelectedItem_t {
	QuoteEnsureVars(&vars)
	var ids []int
	for choiceId := range vars.choices {
		if int(choiceId) <= 0 { continue }
		ids = append(ids, int(choiceId))
	}
	sort.Ints(ids)

	var out []QuoteSelectedItem_t
	for _, itemId := range ids {
		choice := vars.choices[ChoiceId_t(itemId)]
		if int(choice.plan) <= 0 { continue }
		cats := make(map[CategId_t]AddonId_t, len(choice.addons))
		for categId, addon := range choice.addons { cats[categId] = addon }
		out = append(out, QuoteSelectedItem_t{
			itemId: itemId,
			planId: int(choice.plan),
			cats: cats,
		})
	}
	return out
}

func QuoteCloneState(in State_t) State_t {
	out := in
	out.quote = CloneQuoteVars(in.quote)
	return out
}

func QuoteSelectedPreexByCateg(vars QuoteVars_t, itemId int, row QuotePlan_t) map[CategId_t]EuroCent_t {
	out := make(map[CategId_t]EuroCent_t)
	if itemId <= 0 { return out }
	choice, ok := vars.choices[ChoiceId_t(itemId)]
	if !ok { return out }
	for _, preex := range choice.preex {
		base := row.planBase
		if preex.categ > 0 {
			base = 0
			for _, addon := range row.addons {
				if addon.categId != preex.categ { continue }
				base = addon.base
				break
			}
		}
		applied := EuroCent_t(0)
		if preex.amount.euro > 0 { applied = preex.amount.euro }
		if preex.amount.percent > 0 && base > 0 {
			applied = EuroCent_t((int64(base) * int64(preex.amount.percent)) / 10000)
		}
		if applied <= 0 { continue }
		out[preex.categ] += applied
	}
	return out
}

func QuoteSelectedPlanRow(state State_t, item QuoteSelectedItem_t) (QuotePlan_t, bool) {
	work := QuoteCloneState(state)
	QuoteEnsureVars(&work.quote)
	for catId, addon := range item.cats {
		work.quote.planCats[PlanCateg_t{ plan:PlanId_t(item.planId), categ:catId }] = addon
	}

	list := QuotePlans(work).plans
	for _, row := range list {
		if row.planId != item.planId { continue }
		plan, ok := App.lookup.plans.byId[item.planId]
		if ok {
			preexByCateg := QuoteSelectedPreexByCateg(work.quote, item.itemId, row)
			row.commission = QuoteCommission(row, plan, preexByCateg)
		}
		return row, true
	}
	return QuotePlan_t{}, false
}

func QuoteSelectedRows(state State_t) []QuoteSelectedRow_t {
	selected := QuoteSelectedItems(state.quote)
	var out []QuoteSelectedRow_t
	for _, item := range selected {
		row, ok := QuoteSelectedPlanRow(state, item)
		if !ok { continue }
		out = append(out, QuoteSelectedRow_t{ item:item, row:row })
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].row.price != out[j].row.price { return out[i].row.price < out[j].row.price }
		if out[i].row.label != out[j].row.label { return Lower(out[i].row.label) < Lower(out[j].row.label) }
		return out[i].item.itemId < out[j].item.itemId
	})
	return out
}

func QuoteSelectedDrop(state *State_t, itemId int) {
	if itemId <= 0 { return }
	QuoteEnsureDefaults(state)
	delete(state.quote.choices, ChoiceId_t(itemId))
}

func QuoteSelectedAdd(state *State_t, planId int) {
	if planId <= 0 { return }
	QuoteEnsureDefaults(state)
	if len(QuoteSelectedItems(state.quote)) >= quoteSelectedMaxCount { return }

	itemId := int(QuoteAllocChoiceId(&state.quote))
	choice := PlanQuoteInfo_t{
		plan: PlanId_t(planId),
		addons: make(map[CategId_t]AddonId_t),
	}

	rows := QuotePlans(*state).plans
	for _, row := range rows {
		if row.planId != planId { continue }
		for _, addon := range row.addons {
			if !addon.hasMulti || len(addon.choices) == 0 { continue }
			choice.addons[addon.categId] = addon.addon
		}
		break
	}
	state.quote.choices[ChoiceId_t(itemId)] = choice
}

func QuoteSelectedApply(state *State_t, name, value string) bool {
	if planId, ok := QuoteSelectedAddControl(name); ok {
		QuoteSelectedAdd(state, planId)
		return true
	}

	if itemId, ok := QuoteSelectedDelControl(name); ok {
		QuoteSelectedDrop(state, itemId)
		return true
	}

	itemId, catId, ok := QuoteSelectedCatControl(name)
	if !ok { return false }
	QuoteEnsureDefaults(state)
	choice, has := state.quote.choices[ChoiceId_t(itemId)]
	if !has || int(choice.plan) <= 0 { return true }
	if choice.addons == nil { choice.addons = make(map[CategId_t]AddonId_t) }
	choice.addons[catId] = AddonId_t(Atoi(value))
	state.quote.choices[ChoiceId_t(itemId)] = choice
	return true
}
