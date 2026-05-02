package main

import "testing"

func TestEditQSetPreex_ModePersistsOnNoteOnly(t *testing.T) {
	list := []Preex_t{}
	list = EditQSetPreex(list, 0, editQPreexModeEur, ``, ``, true, false, false)
	list = EditQSetPreex(list, 0, ``, ``, `abc`, false, false, true)
	if len(list) != 1 { t.Fatalf("want 1 preex, got %d", len(list)) }

	mode, _, _ := EditQPreexFields(list[0])
	if mode != editQPreexModeEur { t.Fatalf("mode reverted: got %q, want %q", mode, editQPreexModeEur) }
}

func TestEditQSetPreex_AmountThenModeThenNote(t *testing.T) {
	list := []Preex_t{}
	list = EditQSetPreex(list, 0, editQPreexModePct, `10`, ``, false, true, false)
	list = EditQSetPreex(list, 0, editQPreexModeEur, ``, ``, true, false, false)
	list = EditQSetPreex(list, 0, ``, ``, `note`, false, false, true)
	if len(list) != 1 { t.Fatalf("want 1 preex, got %d", len(list)) }

	mode, amount, note := EditQPreexFields(list[0])
	if mode != editQPreexModeEur { t.Fatalf("mode reverted: got %q, want %q", mode, editQPreexModeEur) }
	if amount != `10,00` { t.Fatalf("amount changed: got %q, want 10,00", amount) }
	if note != `note` { t.Fatalf("note changed: got %q, want note", note) }
}

func TestEditQApply_ModePersistsAcrossNoteAndAmount(t *testing.T) {
	state := InitState()
	state.quote = QuoteVars_t{
		choices: map[ChoiceId_t]PlanQuoteInfo_t{
			1: { plan:1, addons:make(map[CategId_t]AddonId_t) },
		},
	}

	if !EditQApply(&state, EditQPreexModeKey(1, 0), editQPreexModeEur) { t.Fatalf("mode apply failed") }
	if !EditQApply(&state, EditQPreexNoteKey(1, 0), `n`) { t.Fatalf("note apply failed") }
	if !EditQApply(&state, EditQPreexAmountKey(1, 0), `10`) { t.Fatalf("amount apply failed") }

	choice := state.quote.choices[1]
	preex, ok := EditQChoicePreex(choice, 0)
	if !ok { t.Fatalf("preex missing") }
	mode, amount, note := EditQPreexFields(preex)
	if mode != editQPreexModeEur { t.Fatalf("mode reverted: got %q, want %q", mode, editQPreexModeEur) }
	if amount != `10,00` { t.Fatalf("amount changed: got %q, want 10,00", amount) }
	if note != `n` { t.Fatalf("note changed: got %q, want n", note) }
}
