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


	xl.WriteQuote()

	status.path = `assets/work`
	err = os.MkdirAll(status.path, 0o775); if err != nil { status.err = err; return status }

	_ = xl.DeleteSheet(xlStyleSheet)

	status.fileName = Str(safeClientName(qvars), ` overview`, If(qvars.slim, `.slim`, ``), `.xlsx`)
	xl.setDefaultCell(quoteSheet, nameCell)
	err = xl.SaveAs(Str(status.path, `/`, status.fileName)); if err != nil { status.err = err; return status }

	return status
}

func (xl *Excel_t)WriteQuote() {
	selected := QuoteSelectedItems(xl.qvars)
	lastBenefitRow := xl.WriteBenefitNames()
	xl.WriteTipsTitle(lastBenefitRow)

	state := InitState()
	state.quote = xl.qvars
	QuoteEnsureVars(&state.quote)

	used := 0
	for _, item := range selected {
		if used >= len(planColumns) { break }
		col := planColumns[used]
		used++

		row, ok := QuoteSelectedPlanRow(state, item)
		if !ok { continue }
		plan, ok := App.lookup.plans.byId[item.planId]
		if !ok { continue }

		xl.WritePlanInfo(lastBenefitRow, col, item, row, plan)
	}

	if used > 0 && used < len(planColumns) {
		lastShown, e0 := sky.ColumnNameToNumber(planColumns[used-1])
		end, e2 := sky.ColumnNameToNumber(planColumns[len(planColumns)-1])
		start := lastShown + 1
		if e0 == nil && e2 == nil && start > 0 && end >= start {
			for col := end; col >= start; col-- {
				colName, e := sky.ColumnNumberToName(col)
				if e != nil { continue }
				_ = xl.RemoveCol(quoteSheet, colName)
			}
		}
	}

	xl.WriteClientInfo()
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
	if len(QuoteSelectedItems(qvars)) == 0 {
		Log(`download-excel blocked: no selected plans`)
		http.Redirect(w, req, `/`, http.StatusSeeOther)
		return
	}

	x := CreateExcelQuote(qvars)
	if x.err != nil {
		Log(x.err)
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, x.path, x.fileName, xlsx)
}
