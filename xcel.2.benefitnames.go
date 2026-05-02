package main

import (
	. "klpm/lib/output"
)

func (xl *Excel_t)WriteBenefitNames() (lastRow int) {
	row := 24
	lastRow = row - 1
	lang := int(xl.qvars.lang)
	deductibleCell := ``

	for _, sec := range App.lookup.benSecs.All() {
		if sec.lang != lang { continue }

		items := xl.benefitItems(sec.section, lang, xl.qvars.slim)
		top, low := row, row+len(items)-1
		aTop, aLow := Str(`A`, top), Str(`A`, low)
		_ = xl.SetCellValue(quoteSheet, aTop, sec.label)
		_ = xl.MergeCell(quoteSheet, aTop, aLow)

		titleStyle, labelStyle := xl.sectionStyles(sec.section)
		_ = xl.SetCellStyle(quoteSheet, aTop, aLow, titleStyle)

		for _, item := range items {
			b := Str(`B`, row)
			_ = xl.SetCellValue(quoteSheet, b, item.label)
			_ = xl.SetCellStyle(quoteSheet, b, b, labelStyle)
			if item.benefit == 1 { deductibleCell = b }
			lastRow = row
			row++
		}
		row++
	}

	if deductibleCell != `` { _ = xl.SetCellStyle(quoteSheet, deductibleCell, deductibleCell, xl.styles[`Spec1Label`]) }

	return lastRow
}

func (xl *Excel_t)sectionStyles(section int) (title, label int) {
	switch section {
	case 0: return xl.styles[`Sec0Title`], xl.styles[`Sec0Label`]
	case 1: return xl.styles[`Sec1Title`], xl.styles[`Sec1Label`]
	case 2: return xl.styles[`Sec2Title`], xl.styles[`Sec2Label`]
	case 3: return xl.styles[`Sec3Title`], xl.styles[`Sec3Label`]
	case 4: return xl.styles[`Sec4Title`], xl.styles[`Sec4Label`]
	}
	return
}

func (xl *Excel_t)benefitItems(section, lang int, slim bool) (out []BenSecItem_t) {
	for _, item := range App.lookup.benSecItems.All() {
		if item.section != section || item.lang != lang { continue }
		if slim && !item.isSlim { continue }
		out = append(out, item)
	}
	return out
}
