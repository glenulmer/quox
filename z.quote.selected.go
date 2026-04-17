package main

import (
	"fmt"
	"sort"

	. "quo2/lib/output"
)

const quoteSelectedSeqKey = `sel-seq`
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

func QuoteSelectedPlanKey(itemId int) string {
	return Str(`selplan-`, itemId)
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

func QuoteSelectedItems(vars UIBagVars_t) []QuoteSelectedItem_t {
	all := make(map[int]QuoteSelectedItem_t)

	for key, value := range vars {
		if itemId, ok := QuoteSelectedPlanControl(key); ok {
			x := all[itemId]
			x.itemId = itemId
			x.planId = Atoi(value)
			if x.cats == nil { x.cats = make(map[CategId_t]AddonId_t) }
			all[itemId] = x
			continue
		}

		itemId, catId, ok := QuoteSelectedCatControl(key)
		if !ok { continue }
		x := all[itemId]
		x.itemId = itemId
		if x.cats == nil { x.cats = make(map[CategId_t]AddonId_t) }
		x.cats[catId] = AddonId_t(Atoi(value))
		all[itemId] = x
	}

	var out []QuoteSelectedItem_t
	for _, x := range all {
		if x.planId <= 0 { continue }
		if x.cats == nil { x.cats = make(map[CategId_t]AddonId_t) }
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].itemId < out[j].itemId })
	return out
}

func QuoteCloneState(in State_t) State_t {
	out := in
	out.quote = make(UIBagVars_t, len(in.quote))
	for k, v := range in.quote { out.quote[k] = v }
	return out
}

func QuoteStateFromVars(vars UIBagVars_t) State_t {
	out := InitState()
	out.quote = make(UIBagVars_t, len(vars))
	for k, v := range vars { out.quote[k] = v }
	return out
}

func QuoteSelectedPlanRow(state State_t, item QuoteSelectedItem_t) (QuotePlan_t, bool) {
	work := QuoteCloneState(state)
	for catId, addon := range item.cats {
		work.quote[QuotePlanCatControlName(item.planId, catId)] = Str(addon)
	}
	list := QuotePlans(work).plans
	for _, row := range list {
		if row.planId == item.planId { return row, true }
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
	if state.quote == nil { return }

	state.quote[QuoteSelectedPlanKey(itemId)] = ``
	for key := range state.quote {
		otherId, _, ok := QuoteSelectedCatControl(key)
		if !ok || otherId != itemId { continue }
		state.quote[key] = ``
	}
}

func QuoteSelectedAdd(state *State_t, planId int) {
	if planId <= 0 { return }
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	if len(QuoteSelectedItems(state.quote)) >= quoteSelectedMaxCount { return }

	itemId := Atoi(state.quote[quoteSelectedSeqKey]) + 1
	state.quote[quoteSelectedSeqKey] = Str(itemId)
	state.quote[QuoteSelectedPlanKey(itemId)] = Str(planId)

	rows := QuotePlans(*state).plans
	for _, row := range rows {
		if row.planId != planId { continue }
		for _, addon := range row.addons {
			if !addon.hasMulti || len(addon.choices) == 0 { continue }
			state.quote[QuoteSelectedCatKey(itemId, addon.categId)] = Str(addon.addon)
		}
		return
	}
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
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	if Trim(state.quote[QuoteSelectedPlanKey(itemId)]) == `` { return true }

	state.quote[QuoteSelectedCatKey(itemId, catId)] = value
	return true
}
