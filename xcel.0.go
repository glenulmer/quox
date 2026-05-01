package main

import (
	"net/http"
	"os"
	"strings"
	
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

type Excel_t struct {
	*sky.File
	qvars QuoteVars_t
	styles map[string]int
	err error
}

const xlsx = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

type xlStatus_t struct {
	path, fileName string
	err error
}

func CreateExcelQuote(qvars QuoteVars_t) (status xlStatus_t) {
	orig, err := sky.OpenFile(`assets/ExcelQuote.xlsx`)
	if err != nil { status.err = err; return status }
	xl := &Excel_t{orig, qvars, nil, nil}
	xl.styles = xl.loadXlStyles()
	defer xl.Close()


	e := xl.WriteQuote(); if e.Err() { status.err = e; return status }

	status.path = `assets/work`
	err = os.MkdirAll(status.path, 0o775); if err != nil { status.err = err; return status }

	_ = xl.DeleteSheet(xlStyleSheet)

	status.fileName = Str(safeClientName(qvars), ` overview`, If(qvars.slim, `.slim`, ``), `.xlsx`)
	xl.setDefaultCell(quoteSheet, nameCell)
	err = xl.SaveAs(Str(status.path, `/`, status.fileName)); if err != nil { status.err = e; return status }

	return status
}

func (xl *Excel_t)WriteQuote() (e checkErr_t) {
	e = xl.WriteClientInfo(); if e.Err() { return e }

	e = xl.WritePlansTop(); if e.Err() { return e }

	lastBenefitRow, e := xl.WriteBenefitNames(); if e.Err() { return e }

	e = xl.WriteTipsTitle(lastBenefitRow); if e.Err() { return e }

	return checkErr_t{}
}


const quoteSheet = `Sheet1`
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

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	qvars := QuoteVars(&state)

	x := CreateExcelQuote(qvars)
	if x.err != nil {
		Log(x.err)
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, x.path, x.fileName, xlsx)
}
