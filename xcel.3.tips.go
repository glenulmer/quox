package main

import (
	. "klpm/lib/output"
)

func (xl *Excel_t)WriteTipsTitle(lastBenefitRow int) (e checkErr_t) {
	i := int(xl.qvars.lang) - 1
	tipTitle := []string{ `Helpful extra tips`, `Hilfreiche Zusatzhinweise` }[i]
	priceMsg := []string{ `Prices subject to increase in January each year.`, `Preise können sich jedes Jahr im Januar erhöhen.` }[i]

	tipTop := lastBenefitRow + 2
	tipLow := tipTop + 2
	priceRow := tipLow + 2
	atop, blow := Str(`A`, tipTop), Str(`B`, tipLow)
	_ = xl.SetCellValue(quoteSheet, atop, tipTitle)
	_ = xl.MergeCell(quoteSheet, atop, blow)
	_ = xl.SetCellStyle(quoteSheet, atop, blow, xl.styles[`TipTitle`])

	aPrice := Str(`A`, priceRow)
	_ = xl.SetCellValue(quoteSheet, aPrice, priceMsg)
	_ = xl.SetCellStyle(quoteSheet, aPrice, aPrice, xl.styles[`SlimNote`])

	return checkErr_t{}
}
