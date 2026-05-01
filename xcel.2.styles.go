package main

import (
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

const xlStyleSheet = `formats`
const xlStyleFrom = `B2`
const xlStyleTo = `J100`

func LoadXlStyles(ex *sky.File) map[string]int {
	out := make(map[string]int)
	if ex == nil { return out }

	left, top, e1 := sky.CellNameToCoordinates(xlStyleFrom)
	right, low, e2 := sky.CellNameToCoordinates(xlStyleTo)
	if e1 != nil || e2 != nil { return out }

	if left > right { left, right = right, left }
	if top > low { top, low = low, top }

	for row := top; row <= low; row++ {
		for col := left; col <= right; col++ {
			cell, e := sky.CoordinatesToCellName(col, row)
			if e != nil { continue }

			name, e := ex.GetCellValue(xlStyleSheet, cell)
			if e != nil || Trim(name) == `` { continue }

			style, e := ex.GetCellStyle(xlStyleSheet, cell)
			if e != nil || style == 0 { continue }
			out[name] = style
		}
	}

	return out
}
