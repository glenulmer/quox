package xl

import (
	sky "github.com/xuri/excelize/v2"
	_ "image/jpeg"
	_ "image/png"
	. "quo2/lib/output"
	"strconv"
	"strings"
)

func Excel(template string, sheets ...string) (ex *Excel_t) {
	var e error
	ex = &Excel_t{File: sky.NewFile(), Styles: make(map[string]int)}
	ex.File, e = sky.OpenFile(template)
	if e != nil {
		panic(Str(e))
		return nil
	}
	ns := len(sheets)
	if ns < 1 {
		ex.Sheet = ex.File.GetSheetName(0)
	} else {
		ex.Sheet = sheets[ns-1]
		ix, ee := ex.File.GetSheetIndex(ex.Sheet)
		if ee != nil || ix < 0 {
			ex = nil
		}
	}
	return ex
}

type Excel_t struct {
	File   *sky.File
	Sheet  string
	Styles map[string]int
}

func (ex *Excel_t) GetSheetList() []string {
	return ex.File.GetSheetList()
}

func (ex *Excel_t) RenameSheet(toName string) {
	_ = ex.File.SetSheetName(ex.Sheet, toName)
	ex.Sheet = toName
}

func (ex *Excel_t) UseSheet(find string) bool {
	ix, e := ex.File.GetSheetIndex(find)
	if e != nil || ix < 0 {
		return false
	}
	ex.Sheet = find
	return true
}

func (ex *Excel_t) CopySheet(orig, new string) error {
	oix, e := ex.File.GetSheetIndex(orig)
	if e != nil || oix < 0 {
		return Error(Str("Can't find sheet '", orig, "'."))
	}
	nix, e := ex.File.NewSheet(new)
	if e != nil || nix < 0 {
		return Error(Str("Can't create new sheet '", new, "'."))
	}
	e = ex.File.CopySheet(oix, nix)
	if e != nil {
		return e
	}
	ex.Sheet = new
	return nil
}

func (ex *Excel_t) DeleteSheet(name string) { _ = ex.File.DeleteSheet(name) }

func (ex *Excel_t) SetColWidth(col string, width float64) {
	col = strings.ToUpper(strings.TrimSpace(col))
	if col == `` {
		return
	}
	_ = ex.File.SetColWidth(ex.Sheet, col, col, width)
}

func (ex *Excel_t) Save(target string) error {
	if ex == nil {
		return Error("nil Excel_t pointer")
	}
	return ex.File.SaveAs(target)
}

func (ex *Excel_t) GetCell(col string, row int) string {
	val, e := ex.File.GetCellValue(ex.Sheet, Str(col, row))
	if e != nil {
		return ``
	}
	return val
}

func (ex *Excel_t) GetCellCents(col string, row int) int {
	cents, ok := CentsFromText(ex.GetCell(col, row))
	if !ok {
		return 0
	}
	return cents
}

func (ex *Excel_t) SetCell(col string, row int, val interface{}, styles ...int) {
	cell := Str(col, row)
	_ = ex.File.SetCellValue(ex.Sheet, cell, val)
	if n := len(styles); n > 0 {
		ex.StyleCell(cell, styles[n-1])
	}
}

func (ex *Excel_t) StyleCell(cell string, style int) {
	ex.StyleCells(cell, cell, style)
}

func (ex *Excel_t) StyleCells(UL, LR string, styleId int) {
	_ = ex.File.SetCellStyle(ex.Sheet, UL, LR, styleId)
}

func Col(col int) (c string) {
	name, e := sky.ColumnNumberToName(col)
	if e != nil {
		return ``
	}
	return name
}

func (ex *Excel_t) GetCellStyle(col string, row int) int {
	style, e := ex.File.GetCellStyle(ex.Sheet, Str(col, row))
	if e != nil {
		return 0
	}
	return style
}

func cellCoords(cell string) (col, row int, ok bool) {
	cell = strings.TrimSpace(cell)
	if cell == `` {
		return 0, 0, false
	}

	k := 0
	for k < len(cell) {
		ch := cell[k]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
			k++
			continue
		}
		break
	}
	if k == 0 || k == len(cell) {
		return 0, 0, false
	}

	col, ok = columnNumber(cell[:k])
	if !ok {
		return 0, 0, false
	}
	row = Atoi(cell[k:])
	if col < 1 || row < 1 {
		return 0, 0, false
	}
	return col, row, true
}

func columnNumber(col string) (int, bool) {
	ix, e := sky.ColumnNameToNumber(strings.ToUpper(strings.TrimSpace(col)))
	if e != nil || ix < 1 {
		return 0, false
	}
	return ix, true
}

func Cell(col, row int) string { return Str(Col(col), row) }

func (ex *Excel_t) AddStyles(sheet string, upperLeft, lowerRight string) {
	saveSheet := ex.Sheet
	ex.Sheet = sheet

	left, upper, ok1 := cellCoords(upperLeft)
	right, lower, ok2 := cellCoords(lowerRight)
	if !ok1 || !ok2 {
		ex.Sheet = saveSheet
		return
	}

	if left > right {
		left, right = right, left
	}
	if upper > lower {
		upper, lower = lower, upper
	}

	for row := upper; row <= lower; row++ {
		for col := left; col <= right; col++ {
			cell := Cell(col, row)
			val, e := ex.File.GetCellValue(sheet, cell)
			if e != nil || len(val) == 0 {
				continue
			}
			style, e := ex.File.GetCellStyle(sheet, cell)
			if e != nil || style == 0 {
				continue
			}
			ex.Styles[val] = style
		}
	}

	ex.Sheet = saveSheet
}

func allDigits(s string) bool {
	if s == `` {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func atoiOK(s string) (int, bool) {
	if !allDigits(s) {
		return 0, false
	}
	n, e := strconv.Atoi(s)
	if e != nil {
		return 0, false
	}
	return n, true
}

func CentsFromText(in string) (int, bool) {
	s := strings.TrimSpace(in)
	if s == `` {
		return 0, true
	}

	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	if s == `` {
		return 0, true
	}

	for _, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
		case ch == '.' || ch == ',':
		default:
			return 0, false
		}
	}

	lastDot := strings.LastIndex(s, ".")
	lastCom := strings.LastIndex(s, ",")
	lastSep := lastDot
	if lastCom > lastSep {
		lastSep = lastCom
	}

	if lastSep < 0 {
		return atoiOK(s)
	}

	right := s[lastSep+1:]
	useDecimal := len(right) > 0 && len(right) <= 2 && allDigits(right)

	if !useDecimal {
		digits := strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", "")
		return atoiOK(digits)
	}

	left := strings.ReplaceAll(strings.ReplaceAll(s[:lastSep], ".", ""), ",", "")
	if left == `` {
		left = `0`
	}
	L, ok := atoiOK(left)
	if !ok {
		return 0, false
	}

	if len(right) == 1 {
		right += `0`
	}
	R, ok := atoiOK(right)
	if !ok {
		return 0, false
	}

	return (L * 100) + R, true
}
