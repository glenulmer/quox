package main

import (
	"net/http"
	"os"
	"strings"

	sky "github.com/xuri/excelize/v2"

	. "quo2/lib/output"
)

const xlsx = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

const template = `assets/ExcelQuote.xlsx`
const workDir = `assets/work`
const sheet = `Sheet1`
const nameCell = `A3`

func IsSlimXl(req *http.Request) bool {
	if req == nil { return false }

	slim := Lower(Trim(req.FormValue(`slim`)))
	if slim == `1` || slim == `true` { return true }

	download := Lower(Trim(req.FormValue(`DownloadExcel`)))
	if Contains(download, `slim=true`) || Contains(download, `slim=1`) { return true }

	return false
}

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
	return Str(name, ` initial overview`, slimPart, `.xlsx`)
}

func CreateXlQuote(vars QuoteVars_t, slim bool) (path, fname string, ok bool) {
	clientName := ClientName(vars)

	ex, e := sky.OpenFile(template)
	if e != nil {
		Log(e)
		return ``, ``, false
	}
	defer ex.Close()

	styles := LoadXlStyles(ex)
	if e = SetXlStyled(ex, styles, sheet, nameCell, clientName, xlStyleClient); e != nil {
		Log(e)
		return ``, ``, false
	}

	path = workDir
	if e = os.MkdirAll(path, 0o775); e != nil {
		Log(e)
		return ``, ``, false
	}

	fname = XlFileName(clientName, slim)
	if e = ex.SaveAs(Str(path, `/`, fname)); e != nil {
		Log(e)
		return ``, ``, false
	}

	return path, fname, true
}
