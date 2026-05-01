package main

import (
	"net/http"
	"os"
	"strings"
	
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

const xlsx = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

const template = `assets/ExcelQuote.xlsx`
const workDir = `assets/work`

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	vars := QuoteVars(&state)

	path, fileName, ok := CreateXlQuote(vars)
	if !ok {
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, path, fileName, xlsx)
}

func CreateXlQuote(vars QuoteVars_t) (path, fileName string, ok bool) {
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

	_ = ex.DeleteSheet(xlStyleSheet)

	fileName = Str(safeClientName(vars), ` overview`, If(vars.slim, `.slim`, ``), `.xlsx`)
	setXlDefaultCell(ex, quoteSheet, nameCell)
	if e = ex.SaveAs(Str(path, `/`, fileName)); e != nil {
		return ``, ``, false
	}

	return path, fileName, true
}

const quoteSheet = `Quote`
const nameCell = `A3`

func setXlDefaultCell(ex *sky.File, tab, cell string) {
	if ix, e := ex.GetSheetIndex(tab); e == nil && ix >= 0 {
		ex.SetActiveSheet(ix)
	}
	_ = ex.SetPanes(tab, &sky.Panes{
		Freeze: false,
		Split: false,
		TopLeftCell: cell,
		Selection: []sky.Selection{
			{SQRef: cell, ActiveCell: cell},
		},
	})
}


func safeClientName(vars QuoteVars_t) string {
	name := Trim(vars.core.clientName)
	if name == `` { return `Customer` }

	name = Replace(name, `/`, `-`)
	name = Replace(name, `\`, `-`)
	name = Replace(name, string(os.PathSeparator), `-`)
	name = strings.Map(func(r rune) rune {
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
	}, name)
	name = Trim(name)
	if name == `` { return `Customer` }
	return name
}
