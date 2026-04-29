package main

import (
	. "klpm/lib/date"
	. "klpm/lib/dec2"
)

type ChoiceId_t int

type QuoteVars_t struct {
	lang LangId_t
	slim int
	sortBy string
	core CoreVars_t
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
