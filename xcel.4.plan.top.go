package main

import (
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

var planColumns = []string{ `D`, `F`, `H`, `J`, `L` }

func (xl *Excel_t)WritePlanTop(col string, item QuoteSelectedItem_t, row QuotePlan_t, plan Plan_t) (e checkErr_t) {
	topNote := App.lookup.planNotes[TopNote_t{plan:PlanId_t(item.planId), lang:xl.qvars.lang}]
	if topNote != `` {
		c3 := Str(col, 3)
		_ = xl.SetCellValue(quoteSheet, c3, topNote)
		_ = xl.SetCellStyle(quoteSheet, c3, c3, xl.styles[`TopNote`])
	}

	_ = xl.SetCellValue(quoteSheet, Str(col, 4), Str(plan.provName, ` / `, plan.name))
	_ = xl.SetCellValue(quoteSheet, Str(col, 5), Str(`Ref `, row.planId, `/`, EuroFlatFromCent(row.commission).Int64()))
	_ = xl.SetCellStyle(quoteSheet, Str(col, 5), Str(col, 5), xl.styles[`Commission`])

	pvn := EuroCent_t(0)
	if addon, ok := QuotePlanAddonByTag(row, `pvn`); ok { pvn = addon.base + addon.surcharge }
	sick := EuroCent_t(0)
	sickText := `-`
	if addon, ok := QuotePlanAddonByCateg(row, catSick); ok {
		sick = addon.base + addon.surcharge
		sickText = QuoteAddonPickText(addon)
		if sickText == `` { sickText = addon.categ }
	}

	preexByCateg := QuoteSelectedPreexByCateg(xl.qvars, item.itemId, row)
	pX := EuroCent_t(0)
	for _, v := range preexByCateg { pX += v }
	hic := row.price - pvn - sick
	hicText := Str(hic)
	if pX > 0 { hicText = Str(hic, ` + `, pX) }
	_ = xl.SetCellValue(quoteSheet, Str(col, 6), hicText)
	_ = xl.SetCellValue(quoteSheet, Str(col, 7), Str(pvn))
	_ = xl.SetCellValue(quoteSheet, Str(col, 8), Str(sick))
	_ = xl.SetCellValue(quoteSheet, Str(col, 9), sickText)

	yourPay := row.price
	if xl.qvars.core.segment == Employee {
		pay := (row.price - pvn) / 2
		if y, ok := App.lookup.years.byId[xl.qvars.core.buy.Year()]; ok {
			maxShare := y.maxshare.ToEuroCent()
			if pay > maxShare { pay = maxShare }
		}
		pay += pvn / 2
		yourPay -= pay
		if yourPay < 0 { yourPay = 0 }
	}

	depState := InitState()
	depState.quote = xl.qvars
	QuoteEnsureVars(&depState.quote)
	for ix, dep := range xl.qvars.dependants {
		if ix >= 10 { break }
		depState.quote.core.birth = dep.birth
		depState.quote.core.vision = dep.vision
		depRow, ok := QuoteSelectedPlanRow(depState, item)
		depPrice := EuroCent_t(0)
		if ok { depPrice = depRow.price }
		_ = xl.SetCellValue(quoteSheet, Str(col, 11+ix), Str(depPrice))
		yourPay += depPrice
	}

	if xl.qvars.core.segment == Employee { _ = xl.SetCellValue(quoteSheet, Str(col, 21), Str(row.price)) }
	_ = xl.SetCellValue(quoteSheet, Str(col, 22), Str(yourPay))

	return checkErr_t{}
}
