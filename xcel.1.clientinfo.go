package main

import (
	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/date"
	. "klpm/lib/output"
)

func (xl *Excel_t)WriteClientInfo() (e checkErr_t) {
	name := Trim(xl.qvars.core.clientName)
	if name == `` { name = clientNameDefault() }
	_ = xl.SetCellValue(quoteSheet, `A3`, name)

	if style := xl.styles[`Client Name`]; style != 0 {
		_ = xl.SetCellStyle(quoteSheet, `A3`, `A3`, style)
	}

	if Valid(xl.qvars.core.birth) {
		line := clientDOB(xl.qvars.lang, xl.qvars.core.birth)
		_ = xl.SetCellValue(quoteSheet, `A4`, line)
	}

	switch xl.qvars.lang {
	case German:
		_ = xl.SetCellValue(quoteSheet, `A6`, `Krankenversicherung`)
		_ = xl.SetCellValue(quoteSheet, `A7`, `Pflegepflichtversicherung`)
	}

	if xl.qvars.core.sickCover > 0 {
		line := clientCover(xl.qvars.lang, xl.qvars.core.sickCover.OutEuro())
		_ = xl.SetCellValue(quoteSheet, `A8`, line)
	}

	for ix, dep := range xl.qvars.dependants {
		if ix >= 10 { break }
		row := 11 + ix
		cell := Str(`A`, row)

		depName := Trim(dep.name)
		if depName == `` { depName = depNameDefault(xl.qvars.lang, ix+1) }

		line := depMonthlyCost(xl.qvars.lang, depName)
		if Valid(xl.qvars.core.buy) && Valid(dep.birth) {
			age := xl.qvars.core.buy.Year() - dep.birth.Year()
			line = depMonthlyCostAge(xl.qvars.lang, line, age)
		}
		_ = xl.SetCellValue(quoteSheet, cell, line)
	}

	e = xl.writeCostRows(); if e.Err() { return e }

	depN := len(xl.qvars.dependants)
	if depN > 10 { depN = 10 }
	for n := 10 - depN; n > 0; n-- {
		e = checkErr_t{xl.RemoveRow(quoteSheet, 11+depN)}; if e.Err() { return e }
	}

	return checkErr_t{}
}

func (xl *Excel_t)writeCostRows() (e checkErr_t) {
	if xl == nil { return checkErr_t{Error(`nil excel file`)} }
	switch xl.qvars.core.segment {
	case Employee:
		f21 := xl.baseFont(quoteSheet, `A21`)
		f22 := xl.baseFont(quoteSheet, `A22`)
		_ = xl.SetCellRichText(quoteSheet, `A21`, empTotalCost(xl.qvars.lang, f21))
		_ = xl.SetCellRichText(quoteSheet, `A22`, empPays(xl.qvars.lang, f22))
	default:
		_ = xl.SetCellValue(quoteSheet, `A22`, otherPays(xl.qvars.lang))
		e = checkErr_t{xl.RemoveRow(quoteSheet, 21)}; if e.Err() { return e }
	}
	return checkErr_t{}
}

func clientNameDefault() string {
	return `Customer`
}

func clientDOB(lang LangId_t, birth CalDate_t) string {
	switch lang { case German: return Str(`Geburtsdatum: `, birth.Format(`d. mon, yyyy`)) }
	return Str(`Date of birth: `, birth.Format(`d. mon, yyyy`))
}

func clientCover(lang LangId_t, cover string) string {
	switch lang { case German: return Str(`Krankentagegeld (bei einem Einkommen von `, cover, `)`) }
	return Str(`Daily Sick Pay (for a `, cover, ` income)`)
}

func depNameDefault(lang LangId_t, n int) string {
	switch lang { case German: return Str(`Angehörige `, n) }
	return Str(`Dependant `, n)
}

func depMonthlyCost(lang LangId_t, depName string) string {
	switch lang { case German: return Str(depName, ` monatlicher Beitrag`) }
	return Str(depName, `'s monthly cost`)
}

func depMonthlyCostAge(lang LangId_t, line string, age int) string {
	switch lang { case German: return Str(line, ` (Alter `, age, `)`) }
	return Str(line, ` (age `, age, `)`)
}

func empTotalCost(lang LangId_t, base sky.Font) []sky.RichTextRun {
	bold, plain := base, base
	bold.Bold, plain.Bold = true, false
	switch lang {
	case German:
		return []sky.RichTextRun{
			{Text:`Monatliche Gesamtkosten `, Font:&bold},
			{Text:`(inkl. Arbeitgeberzuschuss)`, Font:&plain},
		}
	}
	return []sky.RichTextRun{
		{Text:`Total monthly cost `, Font:&bold},
		{Text:`(incl. employer subsidy)`, Font:&plain},
	}
}

func empPays(lang LangId_t, base sky.Font) []sky.RichTextRun {
	bold, sub := base, base
	bold.Bold = true
	sub.Bold = false
	sub.Size = 14
	switch lang {
	case German:
		return []sky.RichTextRun{
			{Text:`Ihr monatlicher Beitrag `, Font:&bold},
			{Text:`(nach Zuschuss)`, Font:&sub},
		}
	}
	return []sky.RichTextRun{
		{Text:`Your monthly cost `, Font:&bold},
		{Text:`(after subsidy)`, Font:&sub},
	}
}

func otherPays(lang LangId_t) string {
	switch lang { case German: return `Ihr monatlicher Beitrag` }
	return `Your monthly cost`
}

func (xl *Excel_t)baseFont(tab, cell string) (out sky.Font) {
	if xl == nil { return out }
	styleId, e := xl.GetCellStyle(tab, cell)
	if e != nil || styleId == 0 { return out }
	style, e := xl.GetStyle(styleId)
	if e != nil || style == nil || style.Font == nil { return out }
	return *style.Font
}
