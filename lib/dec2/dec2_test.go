package dec2

import "testing"

func TestEuroCentFormat(t *testing.T) {
	z := EuroCent_t(56742)
	if got, want := z.String(), "567,42"; got != want {
		t.Fatalf("EuroCent String: got %q want %q", got, want)
	}
	if got, want := z.OutEuro(), "567,42 €"; got != want {
		t.Fatalf("EuroCent OutEuro: got %q want %q", got, want)
	}
}

func TestEuroCentZeroRender(t *testing.T) {
	z := EuroCent_t(0)
	if got, want := z.OutEuro(), "-"; got != want {
		t.Fatalf("EuroCent OutEuro zero: got %q want %q", got, want)
	}
	if got, want := z.HTMLText(), "-"; got != want {
		t.Fatalf("EuroCent HTMLText zero: got %q want %q", got, want)
	}
}

func TestEuroFlatFormat(t *testing.T) {
	z := EuroFlat_t(1889)
	if got, want := z.String(), "1.889"; got != want {
		t.Fatalf("EuroFlat String: got %q want %q", got, want)
	}
	if got, want := z.OutEuro(), "1.889 €"; got != want {
		t.Fatalf("EuroFlat OutEuro: got %q want %q", got, want)
	}
}

func TestCommissionRoundHalfUp(t *testing.T) {
	euro := EuroCent_t(56742)
	months := Months_t(333)

	got := Commission(euro, months)
	if want := EuroCent_t(188951); got != want {
		t.Fatalf("Commission: got %d want %d", got, want)
	}
}

func TestApplyPercent(t *testing.T) {
	amount := EuroCent_t(188951)
	pct := Percent_t(30)

	got := ApplyPercent(amount, pct)
	if want := EuroCent_t(56685); got != want {
		t.Fatalf("ApplyPercent: got %d want %d", got, want)
	}
}

func TestParseValues(t *testing.T) {
	e, err := ParseEuroCent("1.234,56 €")
	if err != nil || e != 123456 {
		t.Fatalf("ParseEuroCent: got (%d, %v)", e, err)
	}

	m, err := ParseMonths("3,33")
	if err != nil || m != 333 {
		t.Fatalf("ParseMonths: got (%d, %v)", m, err)
	}

	p, err := ParsePercent("99%")
	if err != nil || p != 99 {
		t.Fatalf("ParsePercent: got (%d, %v)", p, err)
	}

	f, err := ParseEuroFlat("1.889 €")
	if err != nil || f != 1889 {
		t.Fatalf("ParseEuroFlat: got (%d, %v)", f, err)
	}
}
