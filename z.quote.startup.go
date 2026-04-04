package main

import "strings"

const forcedPlanA, forcedPlanB, forcedPlanC = 61, 175, 189

var forcedQuoteDefaults bool

func SetForcedQuoteDefaults() { forcedQuoteDefaults = true }

func QuoteApplyForcedQuoteDefaults(vars QuoteVars_t) {
	if !forcedQuoteDefaults || vars == nil { return }

	vars[`clientName`] = `Jill Jones`
	vars[`birth`] = `1994-06-15`
	vars[`vision`] = `1`

	QuoteDropKeysByPrefix(vars, `selplan-`)
	QuoteDropKeysByPrefix(vars, `selcat-`)
	vars[quoteSelectedSeqKey] = `0`
	state := State_t{ quote: vars }
	QuoteSelectedAdd(&state, forcedPlanA)
	QuoteSelectedAdd(&state, forcedPlanB)
	QuoteSelectedAdd(&state, forcedPlanC)

	QuoteDropKeysByPrefix(vars, `editq-pre-`)
	QuoteDropKeysByPrefix(vars, `editq-prime-`)
	QuoteDropKeysByPrefix(vars, `editq-dep-`)

	vars[editQPreSeqKey] = `1`
	vars[EditQPreKey(1)] = `Diabetes`

	selected := QuoteSelectedItems(vars)
	for _, item := range selected {
		vars[EditQPrimeModeKey(item.itemId, 0)] = editQPrimeModePct
		vars[EditQPrimeAmountKey(item.itemId, 0)] = `15`
		vars[EditQPrimeNoteKey(item.itemId, 0)] = `Diabetes`
	}

	vars[editQDepSeqKey] = `2`
	vars[EditQDepNameKey(1)] = `Bob`
	vars[EditQDepBirthKey(1)] = `1990-06-15`
	vars[EditQDepVisionKey(1)] = ``
	vars[EditQDepCondSeqKey(1)] = `0`

	vars[EditQDepNameKey(2)] = `Jane`
	vars[EditQDepBirthKey(2)] = QuoteJaneBirthDate()
	vars[EditQDepVisionKey(2)] = `1`
	vars[EditQDepCondSeqKey(2)] = `0`
}

func QuoteJaneBirthDate() string {
	birth := EditQDefaultDependentBirth()
	if len(birth) != len(`2006-01-02`) { return birth }
	return birth[:5] + `07-15`
}

func QuoteDropKeysByPrefix(vars QuoteVars_t, prefix string) {
	for key := range vars {
		if strings.HasPrefix(key, prefix) { delete(vars, key) }
	}
}
