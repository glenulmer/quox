package main

import (
	"net/http"
	"os"
	"strings"
	
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

type Excel_t struct { *sky.File }

const xlsx = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

const template = `assets/ExcelQuote.xlsx`
const workDir = `assets/work`

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	qvars := QuoteVars(&state)

	path, fileName, ok := CreateExcelQuote(qvars)
	if !ok {
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, path, fileName, xlsx)
}

func CreateExcelQuote(qvars QuoteVars_t) (path, fileName string, ok bool) {
	orig, e := sky.OpenFile(template)
	if e != nil {
		Log(e)
		return ``, ``, false
	}
	xl := &Excel_t{orig}
	defer xl.Close()

	styles := xl.loadXlStyles()

	if e = xl.WriteClientInfo(styles, qvars); e != nil {
		Log(e)
		return ``, ``, false
	}

	path = workDir
	if e = os.MkdirAll(path, 0o775); e != nil {
		Log(e)
		return ``, ``, false
	}

	_ = xl.DeleteSheet(xlStyleSheet)

	fileName = Str(safeClientName(qvars), ` overview`, If(qvars.slim, `.slim`, ``), `.xlsx`)
	xl.setDefaultCell(quoteSheet, nameCell)
	if e = xl.SaveAs(Str(path, `/`, fileName)); e != nil {
		return ``, ``, false
	}

	return path, fileName, true
}

const quoteSheet = `Quote`
const nameCell = `A3`

func (xl *Excel_t)setDefaultCell(tab, cell string) {
	if ix, e := xl.GetSheetIndex(tab); e == nil && ix >= 0 {
		xl.SetActiveSheet(ix)
	}
	_ = xl.SetPanes(tab, &sky.Panes{
		Freeze: false,
		Split: false,
		TopLeftCell: cell,
		Selection: []sky.Selection{
			{SQRef: cell, ActiveCell: cell},
		},
	})
}


func safeClientName(qvars QuoteVars_t) string {
	name := Trim(qvars.core.clientName)
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

const xlStyleSheet = `formats`
const xlStyleFrom = `B2`
const xlStyleTo = `E50`

func (xl *Excel_t)loadXlStyles() map[string]int {
	out := make(map[string]int)
	if xl == nil { return out }

	left, top, e1 := sky.CellNameToCoordinates(xlStyleFrom)
	right, low, e2 := sky.CellNameToCoordinates(xlStyleTo)
	if e1 != nil || e2 != nil { return out }

	if left > right { left, right = right, left }
	if top > low { top, low = low, top }

	for row := top; row <= low; row++ {
		for col := left; col <= right; col++ {
			cell, e := sky.CoordinatesToCellName(col, row)
			if e != nil { continue }

			name, e := xl.GetCellValue(xlStyleSheet, cell)
			if e != nil || Trim(name) == `` { continue }

			style, e := xl.GetCellStyle(xlStyleSheet, cell)
			if e != nil || style == 0 { continue }
			out[name] = style
		}
	}

	return out
}
