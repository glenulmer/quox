package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

const xlsx = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

const template = `assets/ExcelQuote.xlsx`
const workDir = `assets/work`

func ClientName(vars QuoteVars_t) string {
	name := Trim(vars.core.clientName)
	if name == `` { return `Customer` }
	return name
}

func SafeClientName(in string) string {
	work := Trim(in)
	if work == `` { return `Customer` }

	work = Replace(work, `/`, `-`)
	work = Replace(work, `\`, `-`)
	work = Replace(work, string(os.PathSeparator), `-`)
	work = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == ' ' || r == '-' || r == '_' || r == '.':
			return r
		}
		return '-'
	}, work)
	work = Trim(work)
	if work == `` { return `Customer` }
	return work
}

func XlFileName(clientName string, slim bool) string {
	name := SafeClientName(clientName)
	slimPart := ``
	if slim { slimPart = `.slim` }
	return Str(name, ` overview`, slimPart, `.xlsx`)
}

func CreateXlQuote(vars QuoteVars_t) (path, fname string, ok bool) {
	slim := vars.slim == 1
	ex, e := sky.OpenFile(template)
	if e != nil {
		Log(e)
		return ``, ``, false
	}
	defer ex.Close()

	styles := LoadXlStyles(ex)
	if e = WriteXlLayout(ex, styles, vars); e != nil {
		Log(e)
		return ``, ``, false
	}

	path = workDir
	if e = os.MkdirAll(path, 0o775); e != nil {
		Log(e)
		return ``, ``, false
	}

	fname = XlFileName(ClientName(vars), slim)
	_ = ex.DeleteSheet(xlStyleSheet)
	SetXlDefaultCell(ex, quoteSheet, nameCell)
	if e = ex.SaveAs(Str(path, `/`, fname)); e != nil {
		Log(e)
		return ``, ``, false
	}

	return path, fname, true
}
