package main

import (
	. "klpm/lib/output"
)

func (xl *Excel_t)WritePlanMain(lastBenefitRow int, col string, row QuotePlan_t, plan Plan_t) (e checkErr_t) {
	lang := int(xl.qvars.lang)
	tipTop := lastBenefitRow + 2

	r := 24
	for _, sec := range App.lookup.benSecs.All() {
		if sec.lang != lang { continue }

		items := xl.benefitItems(sec.section, lang, xl.qvars.slim)
		valueStyle := 0
		switch sec.section {
		case 0: valueStyle = xl.styles[`Sec0Value`]
		case 1: valueStyle = xl.styles[`Sec1Value`]
		case 2: valueStyle = xl.styles[`Sec2Value`]
		case 3: valueStyle = xl.styles[`Sec3Value`]
		case 4: valueStyle = xl.styles[`Sec4Value`]
		}

		for _, ben := range items {
			cell := Str(col, r)
			offer := xl.benefitOffer(row, plan, ben.benefit, lang)
			_ = xl.SetCellValue(quoteSheet, cell, offer)
			_ = xl.SetCellStyle(quoteSheet, cell, cell, valueStyle)
			if ben.benefit == 1 { _ = xl.SetCellStyle(quoteSheet, cell, cell, xl.styles[`Spec1Value`]) }
			r++
		}
		r++
	}

	for ix, tip := range App.lookup.familyTips[FamilyTip_t{family:plan.familyId, lang:xl.qvars.lang}] {
		if ix >= 3 { break }
		cell := Str(col, tipTop+ix)
		_ = xl.SetCellValue(quoteSheet, cell, tip)
		_ = xl.SetCellStyle(quoteSheet, cell, cell, xl.styles[`TipText`])
	}

	return checkErr_t{}
}

func (xl *Excel_t)benefitOffer(row QuotePlan_t, plan Plan_t, benefit, lang int) string {
	switch benefit {
	case 1:
		return Str(row.deduct)
	case 2:
		if Trim(plan.nc.note) == `` { return Str(row.noClaims) }
		return Str(plan.nc.note, row.noClaims)
	}

	offer := App.lookup.bensByFamily[BenFamily_t{benefit:benefit, family:int(plan.familyId), lang:lang}]
	for _, addon := range row.addons {
		if addon.addon == 0 { continue }
		txt := App.lookup.bensByAddon[BenAddon_t{benefit:benefit, addon:int(addon.addon), lang:lang}]
		if Trim(txt) != `` { return txt }
	}
	return offer
}
