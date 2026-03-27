package main

import (
	"net/http"

	. "pm/lib/date"
	. "pm/lib/output"
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
		{ name:`custName`, label:`Customer name`, kind:quoteText, placeholder:`Customer name`, phoneGroup:`top`, desktopGroup:`identity`, phoneSpan:8, desktopSpan:6, defaultValue:QuoteDefaultStatic(``) },
		{ name:`segment`, label:`Segment`, kind:quoteSelect, phoneGroup:`top`, desktopGroup:`identity`, phoneSpan:4, desktopSpan:6, choiceSP:`quo_segments_chooser`, defaultValue:QuoteDefaultSelectFirst(`quo_segments_chooser`) },

		{ name:`birth`, label:`Birth date`, kind:quoteDate, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, defaultValue:func(x QuoteDefaults_t) string { return x.birth.Format(`yyyy-mm-dd`) } },
		{ name:`buy`, label:`Buy date`, kind:quoteDate, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, defaultValue:func(x QuoteDefaults_t) string { return x.buy.Format(`yyyy-mm-dd`) } },
		{ name:`sickCover`, label:`Sick cover`, kind:quoteNumber, phoneGroup:`core`, desktopGroup:`core`, phoneSpan:4, desktopSpan:4, min:0, max:150000, step:1000, defaultValue:QuoteDefaultStatic(`75000`) },

		{ name:`priorCov`, label:`Prior cover`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_priorcov_chooser`, defaultValue:QuoteDefaultSelectFirst(`quo_priorcov_chooser`) },
		{ name:`exam`, label:`Exam`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_noexam_chooser`, defaultValue:QuoteDefaultSelectFirst(`quo_noexam_chooser`) },
		{ name:`specref`, label:`Specialist`, kind:quoteSelect, phoneGroup:`core`, desktopGroup:`filters`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_specialist_chooser`, defaultValue:QuoteDefaultSelectFirst(`quo_specialist_chooser`) },

		{ name:`vision`, label:`Vision`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(`1`) },
		{ name:`tempVisa`, label:`Temp Visa`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },
		{ name:`noPVN`, label:`No PVN`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },
		{ name:`naturalMed`, label:`Natural Med`, kind:quoteCheckbox, phoneGroup:`flags`, desktopGroup:`flags`, phoneSpan:3, desktopSpan:3, defaultValue:QuoteDefaultStatic(``) },

		{ name:`deductibleMin`, label:`Deductible min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_deductibles_chooser`, choiceArgs:QuoteChoiceArgs(adult, !max), defaultValue:QuoteDefaultSelectFirst(`quo_deductibles_chooser`, adult, !max) },
		{ name:`hospitalMin`, label:`Hospital min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catHospital, !max), defaultValue:QuoteDefaultSelectFirst(`quo_level_chooser_max`, catHospital, !max) },
		{ name:`dentalMin`, label:`Dental min`, kind:quoteSelect, phoneGroup:`min`, desktopGroup:`min`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catDental, !max), defaultValue:QuoteDefaultSelectFirst(`quo_level_chooser_max`, catDental, !max) },

		{ name:`deductibleMax`, label:`Deductible max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_deductibles_chooser`, choiceArgs:QuoteChoiceArgs(adult, max), defaultValue:QuoteDefaultSelectFirst(`quo_deductibles_chooser`, adult, max) },
		{ name:`hospitalMax`, label:`Hospital max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catHospital, max), defaultValue:QuoteDefaultSelectFirst(`quo_level_chooser_max`, catHospital, max) },
		{ name:`dentalMax`, label:`Dental max`, kind:quoteSelect, phoneGroup:`max`, desktopGroup:`max`, phoneSpan:4, desktopSpan:4, choiceSP:`quo_level_chooser_max`, choiceArgs:QuoteChoiceArgs(catDental, max), defaultValue:QuoteDefaultSelectFirst(`quo_level_chooser_max`, catDental, max) },
	}
}

func QuoteFieldDefs() []QuoteField_t {
	var out []QuoteField_t
	for _, x := range QuoteControlDefs() {
		out = append(out, QuoteField_t{ name:x.name, label:x.label })
	}
	return out
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

func QuoteControlsByGroup(layout, group string) []QuoteControl_t {
	var out []QuoteControl_t
	for _, x := range QuoteControlDefs() {
		if QuoteControlGroup(x, layout) != group { continue }
		out = append(out, x)
	}
	return out
}

func QuoteVars(state State_t) QuoteVars_t {
	out := QuoteDefaultVars()
	for k, v := range state.quote { out[k] = v }
	out[`sortBy`] = QuoteSortMode(out[`sortBy`])
	return out
}

func QuoteFieldList(vars QuoteVars_t) []QuoteField_t {
	var list []QuoteField_t
	for _, x := range QuoteFieldDefs() {
		x.value = vars[x.name]
		list = append(list, x)
	}
	return list
}

func QuoteAllowsField(name string) bool {
	if name == `sortBy` { return true }
	if _, ok := QuoteControlByName(name); ok { return true }
	_, _, ok := QuotePlanCatControl(name)
	return ok
}

func QuoteDefaultVars() QuoteVars_t {
	ctx := QuoteDefaults()
	out := make(QuoteVars_t)
	for _, x := range QuoteControlDefs() {
		if x.defaultValue == nil {
			out[x.name] = ``
			continue
		}
		out[x.name] = x.defaultValue(ctx)
	}
	out[`sortBy`] = sortByPrice
	return out
}

func QuoteSortMode(v string) string {
	if Lower(Trim(v)) == sortByName { return sortByName }
	return sortByPrice
}

func QuoteApply(state *State_t, name, value string) {
	if !QuoteAllowsField(name) { return }
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	state.quote[name] = value
}

func QuoteApplyForm(state *State_t, req *http.Request) {
	if state.quote == nil { state.quote = QuoteDefaultVars() }
	for _, x := range QuoteControlDefs() {
		state.quote[x.name] = req.FormValue(x.name)
	}
}

func CurrentDBDate() CalDate_t {
	var ymd int
	App.DB.CallRow(`quo_today_get`).Scan(&ymd)
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
	v := state.quote[key]
	if v == `` { v = state.quote[Q(key)] }
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
