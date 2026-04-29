package main

import (
	"time"

	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

const catHospital, catDental = 3, 4
const catSick CategId_t = 1
const specrefAddon = 2
const segmentStudent = 4
const sortByName, sortByPrice = `name`, `price`

type QuoteField_t struct {
	name string
	label string
	value string
}

type QuoteChoice_t struct {
	id int
	label string
}

type QuoteDefaults_t struct {
	today CalDate_t
	birth CalDate_t
	buy CalDate_t
}

type QuoteControl_t struct {
	name string
	label string
	kind string
	placeholder string
	phoneGroup string
	desktopGroup string
	phoneSpan int
	desktopSpan int
	choiceSP string
	choiceArgs func() []any
	defaultValue func(QuoteDefaults_t) string
	min int
	max int
	step int
}

const quoteText, quoteDate, quoteNumber, quoteSelect, quoteCheckbox = `text`, `date`, `number`, `select`, `checkbox`

func QuoteDefaults() QuoteDefaults_t {
	today := CurrentDBDate()
	return QuoteDefaults_t{
		today: today,
		birth: DateFromYMD(today.Year()-32, 6, 15),
		buy: today.Days(40).ToWorkDay(),
	}
}

func QuoteDateAddMonths(in CalDate_t, months int) CalDate_t {
	if !Valid(in) { return 0 }
	t := time.Date(in.Year(), time.Month(in.Month()), in.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, months, 0)
	return DateFromYMD(t.Year(), int(t.Month()), t.Day())
}

func QuoteHasFutureLookupYear(today CalDate_t) bool {
	if !Valid(today) { return false }
	for year, _ := range App.lookup.years.All() {
		if year > today.Year() { return true }
	}
	return false
}

func QuoteBuyMaxDate(today CalDate_t) CalDate_t {
	if !Valid(today) { return 0 }
	if QuoteHasFutureLookupYear(today) { return QuoteDateAddMonths(today, 6) }
	return DateFromYMD(today.Year(), 12, 31)
}

func QuoteParseBuyDate(value string) CalDate_t {
	value = Trim(value)
	if value == `` { return 0 }
	if out := Parse(`yyyymmdd`, value); Valid(out) { return out }
	return Parse(`yyyy-mm-dd`, value)
}

func QuoteParseBirthDate(value string) CalDate_t {
	value = Trim(value)
	if value == `` { return 0 }
	if out := Parse(`yyyymmdd`, value); Valid(out) { return out }
	return Parse(`yyyy-mm-dd`, value)
}

func QuoteBuyBounds() (minDate, maxDate, defaultDate CalDate_t) {
	today := CurrentDBDate()
	if !Valid(today) { return 0, 0, 0 }
	minDate = today
	maxDate = QuoteBuyMaxDate(today)
	defaultDate = today.Days(40).ToWorkDay()
	if Valid(maxDate) && int(defaultDate) > int(maxDate) { defaultDate = maxDate }
	if int(defaultDate) < int(minDate) { defaultDate = minDate }
	return minDate, maxDate, defaultDate
}

func QuoteNormalizeBuyValue(value string) string {
	minDate, maxDate, defaultDate := QuoteBuyBounds()
	buyDate := QuoteParseBuyDate(value)
	if !Valid(buyDate) { buyDate = defaultDate }
	if Valid(minDate) && int(buyDate) < int(minDate) { buyDate = minDate }
	if Valid(maxDate) && int(buyDate) > int(maxDate) { buyDate = maxDate }
	if !Valid(buyDate) { return `` }
	return buyDate.Format(`yyyymmdd`)
}

func QuoteBirthBounds() (minDate, maxDate, defaultDate CalDate_t) {
	today := CurrentDBDate()
	if !Valid(today) { return 0, 0, 0 }
	minDate = QuoteDateAddMonths(today, -12*75)
	maxDate = today.Days(-1)
	defaultDate = DateFromYMD(today.Year()-32, 6, 15)
	if Valid(minDate) && int(defaultDate) < int(minDate) { defaultDate = minDate }
	if Valid(maxDate) && int(defaultDate) > int(maxDate) { defaultDate = maxDate }
	return minDate, maxDate, defaultDate
}

func QuoteNormalizeBirthValue(value string) string {
	minDate, maxDate, defaultDate := QuoteBirthBounds()
	birthDate := QuoteParseBirthDate(value)
	if !Valid(birthDate) { birthDate = defaultDate }
	if Valid(minDate) && int(birthDate) < int(minDate) { birthDate = minDate }
	if Valid(maxDate) && int(birthDate) > int(maxDate) { birthDate = maxDate }
	if !Valid(birthDate) { return `` }
	return birthDate.Format(`yyyymmdd`)
}

func QuoteSickCoverMaxByBuyValue(buyValue string) int {
	year := CurrentDBDate().Year()
	buyDate := QuoteParseBuyDate(buyValue)
	if Valid(buyDate) { year = buyDate.Year() }
	if x, ok := App.lookup.years.byId[year]; ok {
		return int(x.maxCover())
	}
	return 150000
}

func QuoteNormalizeSickCoverValue(value, buyValue string) string {
	v := OnlyDigits(value)
	if v < 0 { v = 0 }
	max := QuoteSickCoverMaxByBuyValue(buyValue)
	if max > 0 && v > max { v = max }
	return Str(v)
}

func QuoteChoiceArgs(args ...any) func() []any {
	return func() []any { return args }
}

func QuoteDefaultStatic(v string) func(QuoteDefaults_t) string {
	return func(QuoteDefaults_t) string { return v }
}

func QuoteDefaultSelectFirst(sp string, args ...any) func(QuoteDefaults_t) string {
	return func(QuoteDefaults_t) string {
		return QuoteChoiceFirst(sp, args...)
	}
}

func QuoteControlDefs() []QuoteControl_t {
	adult := true
	max := true
	return []QuoteControl_t{
		{ name:`clientName`, label:`Client name`, kind:quoteText, placeholder:`Client name`, phoneGroup:`top`, desktopGroup:`identity`, phoneSpan:8, desktopSpan:6, defaultValue:QuoteDefaultStatic(``) },
		{ name:`segment`, label:`Segment`, kind:quoteSelect, phoneGroup:`top`, desktopGroup:`identity`, phoneSpan:4, desktopSpan:6, choiceSP:`klpm_segments_chooser`, defaultValue:QuoteDefaultSelectFirst(`klpm_segments_chooser`) },

		{ name:`birth`, label:`Birth date`, kind:quoteDate, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, defaultValue:func(x QuoteDefaults_t) string { return x.birth.Format(`yyyymmdd`) } },
		{ name:`buy`, label:`Buy date`, kind:quoteDate, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, defaultValue:func(x QuoteDefaults_t) string { return x.buy.Format(`yyyymmdd`) } },
		{ name:`sickCover`, label:`Sick cover`, kind:quoteNumber, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, min:0, max:150000, step:1000, defaultValue:QuoteDefaultStatic(`80000`) },

		{ name:`priorCov`, label:`Prior cover`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_priorcov_chooser`, defaultValue:QuoteDefaultSelectFirst(`klpm_priorcov_chooser`) },
		{ name:`exam`, label:`Exam`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_noexam_chooser`, defaultValue:QuoteDefaultSelectFirst(`klpm_noexam_chooser`) },
		{ name:`specref`, label:`Specialist`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_specialist_chooser`, defaultValue:QuoteDefaultSelectFirst(`klpm_specialist_chooser`) },

		{ name:`vision`, label:`Vision`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(`1`) },
		{ name:`tempVisa`, label:`Temp Visa`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },
		{ name:`noPVN`, label:`No PVN`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },
		{ name:`naturalMed`, label:`Natural Med`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },

		{ name:`deductibleMin`, label:`Deductible min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_deductibles_chooser`, choiceArgs:QuoteChoiceArgs(adult, !max), defaultValue:QuoteDefaultSelectFirst(`klpm_deductibles_chooser`, adult, !max) },
		{ name:`hospitalMin`, label:`Hospital min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catHospital, !max), defaultValue:QuoteDefaultSelectFirst(`klpm_level_chooser_max`, catHospital, !max) },
		{ name:`dentalMin`, label:`Dental min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catDental, !max), defaultValue:QuoteDefaultSelectFirst(`klpm_level_chooser_max`, catDental, !max) },

		{ name:`deductibleMax`, label:`Deductible max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_deductibles_chooser`, choiceArgs:QuoteChoiceArgs(adult, max), defaultValue:QuoteDefaultSelectFirst(`klpm_deductibles_chooser`, adult, max) },
		{ name:`hospitalMax`, label:`Hospital max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catHospital, max), defaultValue:QuoteDefaultSelectFirst(`klpm_level_chooser_max`, catHospital, max) },
		{ name:`dentalMax`, label:`Dental max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`klpm_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catDental, max), defaultValue:QuoteDefaultSelectFirst(`klpm_level_chooser_max`, catDental, max) },
	}
}

func QuoteControlByName(name string) (QuoteControl_t, bool) {
	for _, x := range QuoteControlDefs() {
		if x.name == name { return x, true }
	}
	return QuoteControl_t{}, false
}

func QuoteControlChoices(x QuoteControl_t) []QuoteChoice_t {
	if x.choiceSP == `` { return nil }
	args := []any{}
	if x.choiceArgs != nil { args = x.choiceArgs() }
	return QuoteChoices(x.choiceSP, args...)
}

func QuoteControlGroup(x QuoteControl_t, layout string) string {
	if layout == layoutDesktop { return x.desktopGroup }
	return x.phoneGroup
}

func QuoteControlSpan(x QuoteControl_t, layout string) int {
	if layout == layoutDesktop { return x.desktopSpan }
	return x.phoneSpan
}

func QuoteAllowsField(name string) bool {
	if name == `sortBy` { return true }
	if _, ok := QuoteControlByName(name); ok { return true }
	_, _, ok := QuotePlanCatControl(name)
	return ok
}

func QuoteEnsureVars(vars *QuoteVars_t) {
	if vars == nil { return }
	if vars.lang <= 0 { vars.lang = English }
	vars.sortBy = QuoteSortMode(vars.sortBy)
	if vars.planCats == nil { vars.planCats = make(map[PlanCateg_t]AddonId_t) }
	if vars.choices == nil { vars.choices = make(map[ChoiceId_t]PlanQuoteInfo_t) }
	for choiceId, choice := range vars.choices {
		if choice.addons == nil { choice.addons = make(map[CategId_t]AddonId_t) }
		vars.choices[choiceId] = choice
	}
	for i, dep := range vars.dependants {
		if dep.preexByChoice == nil { dep.preexByChoice = make(map[ChoiceId_t][]Preex_t) }
		if dep.depId <= 0 { dep.depId = i + 1 }
		vars.dependants[i] = dep
	}
}

func QuoteVarsEmpty(vars QuoteVars_t) bool {
	if vars.sortBy != `` { return false }
	if vars.lang > 0 { return false }
	if vars.slim != 0 { return false }
	if Valid(vars.core.buy) || Valid(vars.core.birth) { return false }
	if vars.core.clientName != `` || vars.core.email != `` { return false }
	if vars.core.segment != 0 || vars.core.sickCover != 0 { return false }
	if vars.core.priorCov != 0 || vars.core.exam != 0 || vars.core.specref != 0 { return false }
	if vars.core.vision || vars.core.tempVisa || vars.core.noPVN || vars.core.naturalMed { return false }
	if vars.core.deductible.min != 0 || vars.core.deductible.max != 0 { return false }
	if vars.core.hospital.min != 0 || vars.core.hospital.max != 0 { return false }
	if vars.core.dental.min != 0 || vars.core.dental.max != 0 { return false }
	if len(vars.planCats) > 0 || len(vars.choices) > 0 || len(vars.dependants) > 0 { return false }
	return true
}

func QuoteEnsureDefaults(state *State_t) {
	if state == nil { return }
	if QuoteVarsEmpty(state.quote) {
		state.quote = QuoteDefaultVars()
		return
	}
	QuoteEnsureVars(&state.quote)
}

func QuoteValue(vars QuoteVars_t, name string) string {
	switch name {
	case `clientName`:
		return vars.core.clientName
	case `email`:
		return vars.core.email
	case `segment`:
		return Str(vars.core.segment)
	case `birth`:
		if !Valid(vars.core.birth) { return `` }
		return vars.core.birth.Format(`yyyymmdd`)
	case `buy`:
		if !Valid(vars.core.buy) { return `` }
		return vars.core.buy.Format(`yyyymmdd`)
	case `sickCover`:
		return Str(int(vars.core.sickCover))
	case `priorCov`:
		return Str(vars.core.priorCov)
	case `exam`:
		return Str(vars.core.exam)
	case `specref`:
		return Str(vars.core.specref)
	case `vision`:
		if vars.core.vision { return `1` }
		return ``
	case `tempVisa`:
		if vars.core.tempVisa { return `1` }
		return ``
	case `noPVN`:
		if vars.core.noPVN { return `1` }
		return ``
	case `naturalMed`:
		if vars.core.naturalMed { return `1` }
		return ``
	case `deductibleMin`:
		return Str(int(vars.core.deductible.min))
	case `deductibleMax`:
		return Str(int(vars.core.deductible.max))
	case `hospitalMin`:
		return Str(int(vars.core.hospital.min))
	case `hospitalMax`:
		return Str(int(vars.core.hospital.max))
	case `dentalMin`:
		return Str(int(vars.core.dental.min))
	case `dentalMax`:
		return Str(int(vars.core.dental.max))
	case `sortBy`:
		return QuoteSortMode(vars.sortBy)
	case `lang`:
		return Str(int(vars.lang))
	case `slim`:
		return Str(vars.slim)
	}
	if planId, categId, ok := QuotePlanCatControl(name); ok {
		return Str(int(vars.planCats[PlanCateg_t{ plan:PlanId_t(planId), categ:categId }]))
	}
	if itemId, ok := QuoteSelectedPlanControl(name); ok {
		choice, has := vars.choices[ChoiceId_t(itemId)]
		if !has { return `` }
		return Str(int(choice.plan))
	}
	if itemId, categId, ok := QuoteSelectedCatControl(name); ok {
		choice, has := vars.choices[ChoiceId_t(itemId)]
		if !has { return `` }
		return Str(int(choice.addons[categId]))
	}
	if itemId, categId, ok := EditQPreexModeControl(name); ok {
		choice, has := vars.choices[ChoiceId_t(itemId)]
		if !has { return editQPreexModePct }
		preex, has := EditQChoicePreex(choice, categId)
		if !has { return editQPreexModePct }
		mode, _, _ := EditQPreexFields(preex)
		return mode
	}
	if itemId, categId, ok := EditQPreexAmountControl(name); ok {
		choice, has := vars.choices[ChoiceId_t(itemId)]
		if !has { return `` }
		preex, has := EditQChoicePreex(choice, categId)
		if !has { return `` }
		_, amount, _ := EditQPreexFields(preex)
		return amount
	}
	if itemId, categId, ok := EditQPreexNoteControl(name); ok {
		choice, has := vars.choices[ChoiceId_t(itemId)]
		if !has { return `` }
		preex, has := EditQChoicePreex(choice, categId)
		if !has { return `` }
		_, _, note := EditQPreexFields(preex)
		return note
	}
	if depId, ok := EditQDepNameControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			return dep.name
		}
		return ``
	}
	if depId, ok := EditQDepBirthControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			if !Valid(dep.birth) { return `` }
			return dep.birth.Format(`yyyy-mm-dd`)
		}
		return ``
	}
	if depId, ok := EditQDepVisionControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			if dep.vision { return `1` }
			return ``
		}
		return ``
	}
	if depId, itemId, categId, ok := EditQDepChargeModeControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			preex, has := EditQDepChoicePreex(dep, ChoiceId_t(itemId), categId)
			if !has { return editQPreexModePct }
			mode, _, _ := EditQPreexFields(preex)
			return mode
		}
		return editQPreexModePct
	}
	if depId, itemId, categId, ok := EditQDepChargeAmountControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			preex, has := EditQDepChoicePreex(dep, ChoiceId_t(itemId), categId)
			if !has { return `` }
			_, amount, _ := EditQPreexFields(preex)
			return amount
		}
		return ``
	}
	if depId, itemId, categId, ok := EditQDepChargeNoteControl(name); ok {
		for _, dep := range vars.dependants {
			if dep.depId != depId { continue }
			preex, has := EditQDepChoicePreex(dep, ChoiceId_t(itemId), categId)
			if !has { return `` }
			_, _, note := EditQPreexFields(preex)
			return note
		}
		return ``
	}
	return ``
}

func QuoteSetValue(vars *QuoteVars_t, name, value string) bool {
	if vars == nil { return false }
	QuoteEnsureVars(vars)
	switch name {
	case `clientName`:
		vars.core.clientName = value
		return true
	case `email`:
		vars.core.email = value
		return true
	case `segment`:
		vars.core.segment = Atoi(value)
		return true
	case `birth`:
		vars.core.birth = QuoteParseBirthDate(value)
		return true
	case `buy`:
		vars.core.buy = QuoteParseBuyDate(value)
		return true
	case `sickCover`:
		vars.core.sickCover = EuroFlat_t(Atoi(value))
		return true
	case `priorCov`:
		vars.core.priorCov = Atoi(value)
		return true
	case `exam`:
		vars.core.exam = Atoi(value)
		return true
	case `specref`:
		vars.core.specref = Atoi(value)
		return true
	case `vision`:
		vars.core.vision = QuoteVarBool(value)
		return true
	case `tempVisa`:
		vars.core.tempVisa = QuoteVarBool(value)
		return true
	case `noPVN`:
		vars.core.noPVN = QuoteVarBool(value)
		return true
	case `naturalMed`:
		vars.core.naturalMed = QuoteVarBool(value)
		return true
	case `deductibleMin`:
		vars.core.deductible.min = EuroFlat_t(Atoi(value))
		return true
	case `deductibleMax`:
		vars.core.deductible.max = EuroFlat_t(Atoi(value))
		return true
	case `hospitalMin`:
		vars.core.hospital.min = LevelId_t(Atoi(value))
		return true
	case `hospitalMax`:
		vars.core.hospital.max = LevelId_t(Atoi(value))
		return true
	case `dentalMin`:
		vars.core.dental.min = LevelId_t(Atoi(value))
		return true
	case `dentalMax`:
		vars.core.dental.max = LevelId_t(Atoi(value))
		return true
	case `sortBy`:
		vars.sortBy = QuoteSortMode(value)
		return true
	case `lang`:
		lang := Atoi(value)
		if lang <= 0 { lang = int(English) }
		vars.lang = LangId_t(lang)
		return true
	case `slim`:
		if Atoi(value) == 1 { vars.slim = 1 } else { vars.slim = 0 }
		return true
	}
	if planId, categId, ok := QuotePlanCatControl(name); ok {
		vars.planCats[PlanCateg_t{ plan:PlanId_t(planId), categ:categId }] = AddonId_t(Atoi(value))
		return true
	}
	return false
}

func quoteBaseDefaultVars() QuoteVars_t {
	ctx := QuoteDefaults()
	out := QuoteVars_t{
		lang: English,
		slim: 0,
		sortBy: sortByPrice,
		planCats: make(map[PlanCateg_t]AddonId_t),
		choices: make(map[ChoiceId_t]PlanQuoteInfo_t),
	}
	for _, x := range QuoteControlDefs() {
		value := ``
		if x.defaultValue != nil { value = x.defaultValue(ctx) }
		QuoteSetValue(&out, x.name, value)
	}
	return out
}

func QuoteDefaultQuoteVars() QuoteVars_t {
	if AutoChoose.ready { return CloneQuoteVars(AutoChoose.qvars) }
	return quoteBaseDefaultVars()
}

func QuoteDefaultVars() QuoteVars_t {
	return QuoteDefaultQuoteVars()
}

func QuoteSortMode(v string) string {
	if Lower(Trim(v)) == sortByName { return sortByName }
	return sortByPrice
}

func QuoteApply(state *State_t, name, value string) {
	if !QuoteAllowsField(name) { return }
	QuoteEnsureDefaults(state)
	switch name {
	case `sickCover`:
		value = QuoteNormalizeSickCoverValue(value, QuoteValue(state.quote, `buy`))
	case `birth`:
		value = QuoteNormalizeBirthValue(value)
	case `buy`:
		value = QuoteNormalizeBuyValue(value)
	}
	QuoteSetValue(&state.quote, name, value)
	if name == `buy` {
		cover := QuoteNormalizeSickCoverValue(QuoteValue(state.quote, `sickCover`), value)
		QuoteSetValue(&state.quote, `sickCover`, cover)
	}
}

func CurrentDBDate() CalDate_t {
	var ymd int
	App.DB.CallRow(`klpm_today_get`).Scan(&ymd)
	return CalDate(ymd)
}

func QuoteChoices(sp string, args ...any) []QuoteChoice_t {
	var list []QuoteChoice_t
	rows := App.DB.Call(sp, args...)
	defer rows.Close()
	for rows.Next() {
		var x QuoteChoice_t
		rows.Scan(&x.id, &x.label)
		list = append(list, x)
	}
	return list
}

func QuoteChoiceFirst(sp string, args ...any) string {
	list := QuoteChoices(sp, args...)
	if len(list) == 0 { return `` }
	return Str(list[0].id)
}

func StateValue(state State_t, key string) string {
	v := QuoteValue(state.quote, key)
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
	for _, key := range keys {
		if v, ok := StateIntOK(state, key); ok { return v, true }
	}
	return 0, false
}

func StateDate(state State_t, key string) CalDate_t {
	if key == `birth` { return QuoteParseBirthDate(StateValue(state, key)) }
	if key == `buy` { return QuoteParseBuyDate(StateValue(state, key)) }
	return Parse(`yyyy-mm-dd`, StateValue(state, key))
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
