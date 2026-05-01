package main

import (
	. "klpm/lib/date"
	. "klpm/lib/output"
)

const xlStyleClientName = `Client Name`
const xlStyleYouPay = `YouPay`
const xlSegmentEmployee = 1

type checkErr_t struct { error }
func (e checkErr_t)Err() bool { return e.error != nil }

func (xl *Excel_t)WriteQuote() (e checkErr_t) {
	e = xl.WriteClientInfo(); if e.Err() { return e }
	return checkErr_t{}
}

func (xl *Excel_t)WriteClientInfo() (e checkErr_t) {
	if xl == nil { return checkErr_t{Error(`nil excel file`)} }

	name := Trim(xl.qvars.core.clientName)
	if name == `` { name = `Customer` }
	e = checkErr_t{xl.SetCellValue(quoteSheet, `A3`, name)}; if e.Err() { return e }

	if style := xl.styles[xlStyleClientName]; style != 0 {
		e = checkErr_t{xl.SetCellStyle(quoteSheet, `A3`, `A3`, style)}; if e.Err() { return e }
	}

	if Valid(xl.qvars.core.birth) {
		line := Str(`Date of birth: `, xl.qvars.core.birth.Format(`d. mon, yyyy`))
		e = checkErr_t{xl.SetCellValue(quoteSheet, `A4`, line)}; if e.Err() { return e }
	}

	if xl.qvars.core.sickCover > 0 {
		line := Str(`Daily Sick Pay (for a `, xl.qvars.core.sickCover.OutEuro(), ` income)`)
		e = checkErr_t{xl.SetCellValue(quoteSheet, `A8`, line)}; if e.Err() { return e }
	}

	for ix, dep := range xl.qvars.dependants {
		if ix >= 10 { break }
		row := 11 + ix
		cell := Str(`A`, row)

		depName := Trim(dep.name)
		if depName == `` { depName = Str(`Dependant `, ix+1) }

		line := Str(depName, `'s monthly cost`)
		if Valid(xl.qvars.core.buy) && Valid(dep.birth) {
			age := xl.qvars.core.buy.Year() - dep.birth.Year()
			line = Str(line, ` (age `, age, `)`)
		}
		e = checkErr_t{xl.SetCellValue(quoteSheet, cell, line)}; if e.Err() { return e }
	}

	if xl.qvars.core.segment != xlSegmentEmployee {
		e = checkErr_t{xl.SetCellValue(quoteSheet, `A22`, `Your monthly cost`)}; if e.Err() { return e }
		if style := xl.styles[xlStyleYouPay]; style != 0 {
			e = checkErr_t{xl.SetCellStyle(quoteSheet, `A22`, `A22`, style)}; if e.Err() { return e }
		}
	}

	if xl.qvars.core.segment != xlSegmentEmployee {
		e = checkErr_t{xl.RemoveRow(quoteSheet, 21)}; if e.Err() { return e }
	}

	depN := len(xl.qvars.dependants)
	if depN > 10 { depN = 10 }
	for n := 10 - depN; n > 0; n-- {
		e = checkErr_t{xl.RemoveRow(quoteSheet, 11+depN)}; if e.Err() { return e }
	}

	return checkErr_t{}
}
