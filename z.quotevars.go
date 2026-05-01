package main

import (
	. "klpm/lib/date"
	. "klpm/lib/dec2"
)

type ChoiceId_t int

type QuoteVars_t struct {
	lang LangId_t
	slim bool
	sortBy string
	core CoreVars_t
	planCats map[PlanCateg_t]AddonId_t
	nextChoiceId int
	choices map[ChoiceId_t]PlanQuoteInfo_t
	dependants []Dependant_t
}

type PlanQuoteInfo_t struct {
	plan PlanId_t
	addons map[CategId_t]AddonId_t
	preex []Preex_t
}

type Preex_t struct {
	categ CategId_t
	amount struct { percent Percent_t; euro EuroCent_t }
	note string
}

type Dependant_t struct {
	depId int
	name string
	birth CalDate_t
	vision bool
	preexByChoice map[ChoiceId_t][]Preex_t
}

type CoreVars_t struct {
	clientName string
	email string
	segment int
	birth CalDate_t
	buy CalDate_t
	sickCover EuroFlat_t
	priorCov int
	exam int
	specref int
	vision bool
	tempVisa bool
	noPVN bool
	naturalMed bool
	deductible struct { min, max EuroFlat_t }
	hospital struct { min, max LevelId_t }
	dental struct { min, max LevelId_t }
}

func CloneQuoteVars(in QuoteVars_t) QuoteVars_t {
	out := in
	out.planCats = make(map[PlanCateg_t]AddonId_t, len(in.planCats))
	for k, v := range in.planCats { out.planCats[k] = v }

	out.choices = make(map[ChoiceId_t]PlanQuoteInfo_t, len(in.choices))
	for choiceId, choice := range in.choices {
		next := choice
		next.addons = make(map[CategId_t]AddonId_t, len(choice.addons))
		for categId, addonId := range choice.addons { next.addons[categId] = addonId }
		next.preex = append([]Preex_t{}, choice.preex...)
		out.choices[choiceId] = next
	}

	out.dependants = make([]Dependant_t, 0, len(in.dependants))
	for _, dep := range in.dependants {
		next := dep
		next.preexByChoice = make(map[ChoiceId_t][]Preex_t, len(dep.preexByChoice))
		for choiceId, preex := range dep.preexByChoice {
			next.preexByChoice[choiceId] = append([]Preex_t{}, preex...)
		}
		out.dependants = append(out.dependants, next)
	}
	return out
}
