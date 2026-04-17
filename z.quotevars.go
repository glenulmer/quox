package main

import (
	. "quo2/lib/date"
	. "quo2/lib/dec2"
)

type QuoteVars_t struct {
	core CoreVars_t
	choices []Choice_t
}

type Choice_t struct {
	plan PlanId_t
	products map[CategId_t]ProductId_t
}

type CoreVars_t struct {
	clientName string
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
