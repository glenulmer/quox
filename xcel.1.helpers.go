package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

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
