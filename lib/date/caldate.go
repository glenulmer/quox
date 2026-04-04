package date

import "strconv"
import "strings"
import "time"

type CalDate_t int

const (
	minCalDateYMD = 10000101
	maxCalDateYMD = 99991231
)

var monthShort = [...]string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var monthLong = [...]string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var weekShort = [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var weekLong = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func CalDate(ymd int) CalDate_t {
	if !validCalDateYMD(ymd) { return 0 }
	return CalDate_t(ymd)
}

func DateFromYMD(y, m, d int) CalDate_t { return CalDate(y*10000 + m*100 + d) }

func Valid(in CalDate_t) bool { return in != 0 }

func Parse(format, input string) CalDate_t {
	if len(format) == 0 || len(input) == 0 { return 0 }
	lowerFmt := strings.ToLower(format)
	lowerInput := strings.ToLower(input)
	y, m, d := -1, -1, -1
	weekday := -1
	fi, ii := 0, 0
	for fi < len(format) {
		switch {
		case strings.HasPrefix(lowerFmt[fi:], "yyyy"):
			v, n := parseDigits(input[ii:], 4, 4)
			if n == 0 { return 0 }
			y, fi, ii = v, fi+4, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "mm"):
			v, n := parseDigits(input[ii:], 2, 2)
			if n == 0 { return 0 }
			m, fi, ii = v, fi+2, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "dd"):
			v, n := parseDigits(input[ii:], 2, 2)
			if n == 0 { return 0 }
			d, fi, ii = v, fi+2, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "month"):
			v, n := parseMonthName(lowerInput[ii:], monthLong[:])
			if n == 0 { return 0 }
			m, fi, ii = v, fi+5, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "mon"):
			v, n := parseMonthName(lowerInput[ii:], monthShort[:])
			if n == 0 { return 0 }
			m, fi, ii = v, fi+3, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "m"):
			v, n := parseDigits(input[ii:], 1, 2)
			if n == 0 { return 0 }
			m, fi, ii = v, fi+1, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "weekday"):
			v, n := parseWeekdayName(lowerInput[ii:], weekLong[:])
			if n == 0 { return 0 }
			if weekday >= 0 && weekday != v { return 0 }
			weekday, fi, ii = v, fi+7, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "day"):
			v, n := parseWeekdayName(lowerInput[ii:], weekShort[:])
			if n == 0 { return 0 }
			if weekday >= 0 && weekday != v { return 0 }
			weekday, fi, ii = v, fi+3, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "dth"):
			v, n := parseDigits(input[ii:], 1, 2)
			if n == 0 || ii+n+2 > len(input) { return 0 }
			ii += n
			sfx := lowerInput[ii : ii+2]
			if sfx != daySuffix(v) { return 0 }
			d, fi, ii = v, fi+3, ii+2
		case strings.HasPrefix(lowerFmt[fi:], "yy"):
			v, n := parseDigits(input[ii:], 2, 2)
			if n == 0 { return 0 }
			y, fi, ii = yearFromYY(v), fi+2, ii+n
		case strings.HasPrefix(lowerFmt[fi:], "d"):
			v, n := parseDigits(input[ii:], 1, 2)
			if n == 0 { return 0 }
			d, fi, ii = v, fi+1, ii+n
		default:
			if fi < len(format) && isSep(format[fi]) {
				for fi < len(format) && isSep(format[fi]) {
					fi++
				}
				n := 0
				for ii < len(input) && isSep(input[ii]) {
					ii++
					n++
				}
				if n == 0 { return 0 }
				continue
			}
			if ii >= len(input) || lowerFmt[fi] != lowerInput[ii] { return 0 }
			fi++
			ii++
		}
	}
	for ii < len(input) && isSep(input[ii]) {
		ii++
	}
	if ii != len(input) { return 0 }
	if y < 0 || m < 0 || d < 0 { return 0 }
	if weekday >= 0 && weekdayFromYMD(y, m, d) != weekday { return 0 }
	return CalDate(y*10000 + m*100 + d)
}

func (in CalDate_t)Year() int { return int(in) / 10000 }
func (in CalDate_t)Month() int { return (int(in) / 100) % 100 }
func (in CalDate_t)Day() int { return int(in) % 100 }
func (in CalDate_t)Days(days int) CalDate_t {
	if !Valid(in) { return 0 }
	t := time.Date(in.Year(), time.Month(in.Month()), in.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)
	return CalDate(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

func (in CalDate_t)ToWorkDay() CalDate_t {
	if !Valid(in) { return 0 }
	switch weekdayFromYMD(in.Year(), in.Month(), in.Day()) {
	case int(time.Saturday):
		return in.Days(2)
	case int(time.Sunday):
		return in.Days(1)
	default:
		return in
	}
}

func (in CalDate_t)String() string { return in.Format("yyyymmdd") }

func (in CalDate_t)Hyphens() string { return in.Format("yyyy-mm-dd") }

func (in CalDate_t)Format(format string) string {
	if !Valid(in) { return "00000000" }
	y := in.Year()
	m := in.Month()
	d := in.Day()
	weekday := weekdayFromYMD(y, m, d)
	lower := strings.ToLower(format)
	var out strings.Builder
	out.Grow(len(format) + 8)
	for i := 0; i < len(format); {
		switch {
		case strings.HasPrefix(lower[i:], "yyyy"):
			out.WriteString(pad4(y))
			i += 4
		case strings.HasPrefix(lower[i:], "mm"):
			out.WriteString(pad2(m))
			i += 2
		case strings.HasPrefix(lower[i:], "dd"):
			out.WriteString(pad2(d))
			i += 2
		case strings.HasPrefix(lower[i:], "month"):
			out.WriteString(monthLabel(monthLong[:], m))
			i += 5
		case strings.HasPrefix(lower[i:], "mon"):
			out.WriteString(monthLabel(monthShort[:], m))
			i += 3
		case strings.HasPrefix(lower[i:], "m"):
			out.WriteString(strconv.Itoa(m))
			i++
		case strings.HasPrefix(lower[i:], "weekday"):
			out.WriteString(applyCaseStyle(format[i:i+7], weekLong[weekday]))
			i += 7
		case strings.HasPrefix(lower[i:], "day"):
			out.WriteString(applyCaseStyle(format[i:i+3], weekShort[weekday]))
			i += 3
		case strings.HasPrefix(lower[i:], "dth"):
			out.WriteString(dayWithSuffix(d))
			i += 3
		case strings.HasPrefix(lower[i:], "yy"):
			out.WriteString(pad2(y % 100))
			i += 2
		case strings.HasPrefix(lower[i:], "d"):
			out.WriteString(strconv.Itoa(d))
			i++
		default:
			out.WriteByte(format[i])
			i++
		}
	}
	return out.String()
}

func parseDigits(input string, minDigits, maxDigits int) (value, n int) {
	for n < len(input) && n < maxDigits {
		ch := input[n]
		if ch < '0' || ch > '9' {
			break
		}
		value = value*10 + int(ch-'0')
		n++
	}
	if n < minDigits { return 0, 0 }
	return value, n
}

func parseMonthName(inputLower string, names []string) (month, n int) {
	for i := 1; i <= 12; i++ {
		name := strings.ToLower(names[i])
		if strings.HasPrefix(inputLower, name) {
			return i, len(name)
		}
	}
	return 0, 0
}

func parseWeekdayName(inputLower string, names []string) (weekday, n int) {
	for i := 0; i < 7; i++ {
		name := strings.ToLower(names[i])
		if strings.HasPrefix(inputLower, name) {
			return i, len(name)
		}
	}
	return 0, 0
}

func applyCaseStyle(token, val string) string {
	if token == strings.ToUpper(token) { return strings.ToUpper(val) }
	if token == strings.ToLower(token) { return strings.ToLower(val) }
	return val
}

func isSep(ch byte) bool {
	return ch == '.' || ch == '/' || ch == ' ' || ch == '-'
}

func yearFromYY(yy int) int {
	if yy <= 69 { return 2000 + yy }
	return 1900 + yy
}

func weekdayFromYMD(year, month, day int) int {
	return int(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Weekday())
}

func monthLabel(names []string, month int) string {
	if month < 1 || month > 12 { return "" }
	return names[month]
}

func dayWithSuffix(day int) string {
	if day < 1 { return strconv.Itoa(day) }
	return strconv.Itoa(day) + daySuffix(day)
}

func daySuffix(day int) string {
	if day%100 >= 11 && day%100 <= 13 { return "th" }
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func pad2(in int) string {
	if in >= 0 && in < 10 { return "0" + strconv.Itoa(in) }
	if in > -10 && in < 0 { return "-0" + strconv.Itoa(-in) }
	return strconv.Itoa(in)
}

func pad4(in int) string {
	if in >= 0 && in < 10 { return "000" + strconv.Itoa(in) }
	if in >= 10 && in < 100 { return "00" + strconv.Itoa(in) }
	if in >= 100 && in < 1000 { return "0" + strconv.Itoa(in) }
	if in > -10 && in < 0 { return "-000" + strconv.Itoa(-in) }
	if in <= -10 && in > -100 { return "-00" + strconv.Itoa(-in) }
	if in <= -100 && in > -1000 { return "-0" + strconv.Itoa(-in) }
	return strconv.Itoa(in)
}

func validCalDateYMD(ymd int) bool {
	if ymd < minCalDateYMD || ymd > maxCalDateYMD { return false }
	y := ymd / 10000
	m := (ymd / 100) % 100
	d := ymd % 100
	if m < 1 || m > 12 { return false }
	if d < 1 { return false }
	return d <= daysInMonth(y, m)
}

func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear(year) { return 29 }
		return 28
	}
	return 0
}

func isLeapYear(year int) bool {
	if year%400 == 0 { return true }
	if year%100 == 0 { return false }
	return year%4 == 0
}
