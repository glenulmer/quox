package date

import "testing"

func TestCalDateFormatTokens(t *testing.T) {
	in := CalDate(20260315)
	if got := in.String(); got != "20260315" {
		t.Fatalf("String() = %q, want %q", got, "20260315")
	}
	if got := in.Format("yyyy-mm-dd"); got != "2026-03-15" {
		t.Fatalf("Format yyyy-mm-dd = %q", got)
	}
	if got := in.Format("dth mon yyyy"); got != "15th Mar 2026" {
		t.Fatalf("Format dth mon yyyy = %q", got)
	}
	if got := in.Format("MONTH d, YY"); got != "March 15, 26" {
		t.Fatalf("Format MONTH d, YY = %q", got)
	}
	if got := in.Format("m/d/yy"); got != "3/15/26" {
		t.Fatalf("Format m/d/yy = %q", got)
	}
	if got := in.Format("mm.dd.yyyy"); got != "03.15.2026" {
		t.Fatalf("Format mm.dd.yyyy = %q", got)
	}
	if got := in.Format("Day"); got != "Sun" {
		t.Fatalf("Format Day = %q", got)
	}
	if got := in.Format("DAY"); got != "SUN" {
		t.Fatalf("Format DAY = %q", got)
	}
	if got := in.Format("day"); got != "sun" {
		t.Fatalf("Format day = %q", got)
	}
	if got := in.Format("Weekday"); got != "Sunday" {
		t.Fatalf("Format Weekday = %q", got)
	}
}

func TestCalDateParseTokens(t *testing.T) {
	if got := Parse("yyyy-mm-dd", "2026-03-15"); got != CalDate(20260315) {
		t.Fatalf("Parse yyyy-mm-dd = %v", got)
	}
	if got := Parse("yyyy-mm-dd", "2026/03.15"); got != CalDate(20260315) {
		t.Fatalf("Parse mixed separators = %v", got)
	}
	if got := Parse("dth mon yyyy", "15th Mar 2026"); got != CalDate(20260315) {
		t.Fatalf("Parse dth mon yyyy = %v", got)
	}
	if got := Parse("MONTH d, YY", "march 15, 26"); got != CalDate(20260315) {
		t.Fatalf("Parse MONTH d, YY = %v", got)
	}
	if got := Parse("Weekday yyyy-mm-dd", "Sunday 2026-03-15"); got != CalDate(20260315) {
		t.Fatalf("Parse Weekday yyyy-mm-dd = %v", got)
	}
	if got := Parse("day yyyy-mm-dd", "sun 2026-03-15"); got != CalDate(20260315) {
		t.Fatalf("Parse day yyyy-mm-dd = %v", got)
	}
	if got := Parse("Day yyyy-mm-dd", "Mon 2026-03-15"); got != 0 {
		t.Fatalf("Parse mismatched weekday = %v, want 0", got)
	}
	if got := Parse("yyyy-mm-dd", "2026-02-30"); got != 0 {
		t.Fatalf("Parse invalid date = %v, want 0", got)
	}
	if got := Parse("yyyy-mm-dd", "2026_03_15"); got != 0 {
		t.Fatalf("Parse invalid separator = %v, want 0", got)
	}
	f := "dth mon yyyy"
	in := CalDate(20260315)
	if got := Parse(f, in.Format(f)); got != in {
		t.Fatalf("round-trip parse(format()) = %v, want %v", got, in)
	}
}

func TestCalDateZeroAndInvalid(t *testing.T) {
	if got := CalDate(20260230); got != 0 {
		t.Fatalf("CalDate invalid date = %v, want 0", got)
	}
	var z CalDate_t
	if Valid(z) {
		t.Fatalf("Valid(zero) = true, want false")
	}
	if got := z.String(); got != "00000000" {
		t.Fatalf("zero.String() = %q, want 00000000", got)
	}
	if got := z.Format("yyyy-mm-dd"); got != "00000000" {
		t.Fatalf("zero.Format() = %q, want 00000000", got)
	}
	if got := Parse("yyyy-mm-dd", ""); got != 0 {
		t.Fatalf("Parse empty = %v, want 0", got)
	}
}

func TestCalDateDaysAndToWorkDay(t *testing.T) {
	thu := CalDate(20260312)
	fri := CalDate(20260313)
	sat := CalDate(20260314)
	sun := CalDate(20260315)
	mon := CalDate(20260316)

	if got := fri.Days(3); got != mon {
		t.Fatalf("fri.Days(3) = %v, want %v", got, mon)
	}
	if got := fri.Days(-1); got != thu {
		t.Fatalf("fri.Days(-1) = %v, want %v", got, thu)
	}
	if got := CalDate_t(0).Days(1); got != 0 {
		t.Fatalf("zero.Days(1) = %v, want 0", got)
	}

	if got := fri.ToWorkDay(); got != fri {
		t.Fatalf("fri.ToWorkDay() = %v, want %v", got, fri)
	}
	if got := sat.ToWorkDay(); got != mon {
		t.Fatalf("sat.ToWorkDay() = %v, want %v", got, mon)
	}
	if got := sun.ToWorkDay(); got != mon {
		t.Fatalf("sun.ToWorkDay() = %v, want %v", got, mon)
	}
	if got := CalDate_t(0).ToWorkDay(); got != 0 {
		t.Fatalf("zero.ToWorkDay() = %v, want 0", got)
	}
}
