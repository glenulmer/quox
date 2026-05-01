package main

import (
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

const quoteSheet = `Quote`
const nameCell = `A3`

func SetXlDefaultCell(ex *sky.File, tab, cell string) {
	if ex == nil { return }
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

func WriteXlLayout(ex *sky.File, styles map[string]int, vars QuoteVars_t) error {
	return Error(`Shit code`)
}
