package decimal

import . "klpm/lib/output"

type tVal int64
type tPrec uint8

type tDec struct {
	val  tVal
	prec tPrec
}
type Decimal_t = tDec

func (d tDec)dump() string { return Str("(Dec(", d.val, d.prec, "))") }

var Zero = tDec{}

func Decimal(value tVal, prec tPrec) tDec {
	return tDec{ val: value, prec: prec }
}

func fixPrec(prec tPrec) tPrec {
	if prec > 9 { prec = 9 } else if prec < 0 { prec = 0 }
	return prec
}

func mod(prec tPrec) tVal {
	switch prec {
		case 0: return 1
		case 1: return 10
		case 2: return 100
		case 3: return 1000
		case 4: return 10000
		case 5: return 100000
		case 6: return 1000000
		case 7: return 10000000
		case 8: return 100000000
		case 9: return 1000000000
	}
	var mod = tVal(1)
	for prec > 0 { mod *= 10; prec-- }
	return mod
}

func (p tDec)Val() tVal { return p.val }
func (p tDec)Prec() tPrec { return p.prec }

func (p tDec)SetValue(val tVal) tDec {
	return tDec{ val: val, prec: p.prec }
}

func (p tDec)SetPrec(new tPrec) tDec {
	switch {
		case new > p.prec:
			diff := new - p.prec
			p.val = p.val * mod(diff)
			p.prec = new
		case p.prec > new:
			diff := p.prec - new
			var neg = p.val < 0
			if neg { p.val = 0-p.val }
			p.val = p.val / mod(diff-1)
			round := p.val % 10
			p.val = p.val / 10
			if round > 4 { p.val++ }		
			if neg { p.val = 0-p.val }
			p.prec = new
		}
	return p
}

func (p tDec)Multiply(q tDec) tDec { // limit precision of the result to p.prec
	return tDec { val: p.val * q.val, prec: p.prec + q.prec }.SetPrec(p.prec)
}

func (d tDec)chopt0s() tDec { // chop trailing decimal-place zeroes
	places := d.prec
	for places > 0 && (d.val % 10) == 0 {
		d.val = d.val / 10
		places--
	}
	d.prec = places
	return d
}

func (n tDec)DivideBy(d tDec, places tPrec) tDec {
	if d.val == 0 { panic("Attempt to divide decimal by zero") }
	if n.val == 0 { return tDec{} }

	var flip = false
	if (n.val < 0) != (d.val < 0) { flip = true; n.val = 0-n.val }
	
	n = n.chopt0s()
	d = d.chopt0s()

	x := int(d.prec) - int(n.prec)

	var loop = int(places)
	for x < 0 { x++; loop-- }
	if loop < 0 { return tDec{ prec:places } }
	for x > 0 { x--; n.val *= 10 }

	var val tVal = n.val / d.val
	var rem tVal = n.val % d.val

	if rem > 0 {
		for k := 0; k < loop; k++ {
			rem = rem * 10
			var div tVal = rem / d.val
			val = (val * 10) + div
			rem = rem % d.val
		}
	}

	if rem * 10 / d.val > 4 { val++ }

	if flip { val = 0-val }
	return tDec{ val: val, prec: places }
}

///// PrintOptions stuff

type tPrintOptions struct { dec, thou, pre, post string; usePrec bool; prec tPrec }
func PrintOptions(decimal, thousep string) tPrintOptions { return tPrintOptions{ dec: decimal, thou: thousep } }
func (props tPrintOptions)Decimal() string { return props.dec }
func (props tPrintOptions)Thousands() string { return props.thou }
func (props tPrintOptions)Pre() string { return props.pre }
func (props tPrintOptions)Post() string { return props.post }
func (props tPrintOptions)UsePrec() bool { return props.usePrec }
func (props tPrintOptions)Prec() tPrec { return props.prec }

func (props tPrintOptions)SetPre(s string) tPrintOptions { props.pre = s; return props }
func (props tPrintOptions)SetPost(s string) tPrintOptions { props.post = s; return props }

func (props tPrintOptions)SetPrec(p tPrec) tPrintOptions { props.usePrec = true; props.prec = p; return props }
func (props tPrintOptions)NoPrec() tPrintOptions { props.usePrec = false; props.prec = 0; return props }

var EN = PrintOptions(".", ",")
var DE = PrintOptions(",", ".")
var noPrintOptions = PrintOptions(".", "")

func (d tDec)String() string {
	return noPrintOptions.String(d)
}

func (props tPrintOptions)String(d tDec) string {
	if props.UsePrec() { d = d.SetPrec(props.Prec()) }

	var neg = d.val < 0
	if neg { d.val = 0-d.val }
	var mod = mod(d.prec)
	var left = d.val / mod
	var right = d.val % mod
	var s string
	for left > 999 {
		section := left % 1000
		s = props.Thousands() + Strf("%03d", section) + s
		left = left / 1000
	}
	if neg { left = 0-left }
	s = Str(left) + s
	if right != 0 {
		fmt := Str("%0", d.prec, "d")
		s = s + props.Decimal() + Strf(fmt, right)
	}
	return props.Pre() + s + props.Post()
}

func (props tPrintOptions)DecFromString(s string) tDec {
	var places tPrec
	var digits string

	s = Trim(s) // just because
	s = Replace(s, props.Thousands(), "")
	parts := Split(s, props.Decimal())

	digits = parts[0]

	if len(parts) > 1 {
		right := parts[1]
		places = tPrec(len(right))
		digits += right
	}

	x := tDec{ val: tVal(Atoi(digits)), prec: places }
	if props.UsePrec() { x = x.SetPrec(props.Prec()) }
	return x
}

func DecFromString(s string) tDec { return noPrintOptions.DecFromString(s) }
func DecimalFromString(s string) tDec { return noPrintOptions.DecFromString(s) }

/*
var ECents = PrintOptions(",", ".").SetPrec(2).SetPost(" €")
var ECentsHtml = PrintOptions(",", ".").SetPrec(2).SetPost(" &euro;")
var Euros  = PrintOptions(",", ".").SetPrec(0).SetPost(" €")

func main() {
	x := DecimalFromString("3.1415629")
	Log(DE.String(x))
	Log(ECents.String(x))
	Log(ECentsHtml.String(x))
	Log(Euros.String(x))
}
*/
