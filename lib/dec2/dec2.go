package dec2

import (
	"fmt"
	"strconv"
	"strings"
)

type fix2_t int64

type EuroCent_t fix2_t

type EuroFlat_t fix2_t

type Months_t fix2_t

type Percent_t fix2_t

func (z fix2_t) Int64() int64   { return int64(z) }
func (z fix2_t) String() string { return formatDec2(int64(z)) }

func (z EuroCent_t) Int64() int64 { return int64(z) }
func (z EuroFlat_t) Int64() int64 { return int64(z) }
func (z Months_t) Int64() int64   { return int64(z) }
func (z Percent_t) Int64() int64  { return int64(z) }

func (z EuroCent_t) String() string { return formatDec2(int64(z)) }
func (z EuroFlat_t) String() string { return formatWhole(int64(z)) }
func (z Months_t) String() string   { return formatDec2(int64(z)) }
func (z Percent_t) String() string  { return strconv.FormatInt(int64(z), 10) }

func (z EuroCent_t) OutEuro() string {
	if z == 0 {
		return `-`
	}
	return z.String() + " \u20ac"
}
func (z EuroFlat_t) OutEuro() string { return z.String() + " \u20ac" }

func (z EuroCent_t) HTMLText() string {
	if z == 0 {
		return `-`
	}
	return z.OutEuro()
}
func (z EuroFlat_t) HTMLText() string { return z.OutEuro() }

func (z EuroCent_t) ExcelNumber() float64 { return float64(z) / 100 }
func (z EuroFlat_t) ExcelNumber() float64 { return float64(z) }
func (z Months_t) ExcelNumber() float64   { return float64(z) / 100 }
func (z Percent_t) ExcelNumber() float64  { return float64(z) / 100 }

func (EuroCent_t) ExcelFormat() string { return `#.##0,00 [$€-1]` }
func (EuroFlat_t) ExcelFormat() string { return `#.##0 [$€-1]` }
func (Months_t) ExcelFormat() string   { return `0,00` }
func (Percent_t) ExcelFormat() string  { return `0\%` }

func (z EuroFlat_t) ToEuroCent() EuroCent_t { return EuroCent_t(int64(z) * 100) }

func EuroFlatFromCent(c EuroCent_t) EuroFlat_t {
	if c >= 0 {
		return EuroFlat_t((int64(c) + 50) / 100)
	}
	return EuroFlat_t(-(((-int64(c)) + 50) / 100))
}

func Commission(euro EuroCent_t, months Months_t) EuroCent_t {
	return EuroCent_t(roundHalfUp(int64(euro)*int64(months), 100))
}

func ApplyPercent(amount EuroCent_t, pct Percent_t) EuroCent_t {
	return EuroCent_t(roundHalfUp(int64(amount)*int64(pct), 100))
}

func ParseEuroCent(in string) (EuroCent_t, error) {
	v, e := parseDec2(in)
	return EuroCent_t(v), e
}

func ParseMonths(in string) (Months_t, error) {
	v, e := parseDec2(in)
	return Months_t(v), e
}

func ParsePercent(in string) (Percent_t, error) {
	s := strings.TrimSpace(in)
	s = strings.TrimSuffix(s, "%")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, fmt.Errorf("empty percent")
	}

	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return 0, fmt.Errorf("bad percent: %w", e)
	}
	return Percent_t(v), nil
}

func ParseEuroFlat(in string) (EuroFlat_t, error) {
	s := strings.TrimSpace(in)
	s = strings.TrimSuffix(s, "\u20ac")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, fmt.Errorf("empty euro")
	}

	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return 0, fmt.Errorf("bad euro: %w", e)
	}
	return EuroFlat_t(v), nil
}

func parseDec2(in string) (int64, error) {
	s := strings.TrimSpace(in)
	s = strings.TrimSuffix(s, "\u20ac")
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}

	sign := int64(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = strings.TrimPrefix(s, "-")
	}

	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "")

	whole := s
	dec := "00"
	if strings.Contains(s, ",") {
		p := strings.Split(s, ",")
		if len(p) != 2 {
			return 0, fmt.Errorf("bad decimal format")
		}
		whole = p[0]
		dec = p[1]
		switch len(dec) {
		case 0:
			dec = "00"
		case 1:
			dec = dec + "0"
		case 2:
		default:
			return 0, fmt.Errorf("too many decimal digits")
		}
	}

	if whole == "" {
		whole = "0"
	}
	if !digitsOnly(whole) || !digitsOnly(dec) {
		return 0, fmt.Errorf("non-digit value")
	}

	wi, e1 := strconv.ParseInt(whole, 10, 64)
	if e1 != nil {
		return 0, fmt.Errorf("bad whole part: %w", e1)
	}
	di, e2 := strconv.ParseInt(dec, 10, 64)
	if e2 != nil {
		return 0, fmt.Errorf("bad decimal part: %w", e2)
	}

	return sign * (wi*100 + di), nil
}

func roundHalfUp(numer, denom int64) int64 {
	if denom <= 0 {
		panic("denominator must be positive")
	}
	if numer >= 0 {
		return (numer + (denom / 2)) / denom
	}
	return -(((-numer) + (denom / 2)) / denom)
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func formatDec2(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	whole := v / 100
	dec := v % 100
	return sign + formatWhole(whole) + "," + fmt.Sprintf("%02d", dec)
}

func formatWhole(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}

	s := strconv.FormatInt(v, 10)
	n := len(s)
	if n <= 3 {
		return sign + s
	}

	out := make([]byte, 0, n+(n/3))
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, s[i])
	}
	return sign + string(out)
}
