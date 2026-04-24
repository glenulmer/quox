package main

import (
	"sort"

	sky "github.com/xuri/excelize/v2"

	. "quo2/lib/date"
	. "quo2/lib/dec2"
	. "quo2/lib/output"
)

const sheet = `Sheet1`
const nameCell = `A3`
const birthCell = `A4`
const sickCell = `A8`
const payCell = `A22`
const logoRow = 2
const topNoteRow = 3
const planNameRow = 4
const planRefRow = 5
const hicRow = 6
const pvnRow = 7
const sickRow = 8
const sickAfterRow = 9
const firstDepRow = 11
const maxDeps = 10
const monthWithEmpRow = 21
const monthYouPayRow = 22
const segmentEmployee = 1
const slimRowLow = 6
const slimRowHigh = 11
const studentNoSickRowLow = 8
const studentNoSickRowHigh = 9
const benefitRow = 24
const deductibleBenefit = 1
const noClaimBenefit = 2
var xlPlanCols = []string{`D`, `F`, `H`, `J`, `L`}
const xlPlanColDeletes = 10
const xlLogoDir = `assets/logos/`

func BirthLine(vars QuoteVars_t) string {
	if !Valid(vars.core.birth) { return `` }
	return Str(`Date of birth: `, vars.core.birth.Format(`d. mon, yyyy`))
}

func SickCoverLine(vars QuoteVars_t) string {
	if vars.core.sickCover <= 0 { return `` }
	return Str(`Daily Sick Pay (for a `, vars.core.sickCover.OutEuro(), ` income)`)
}

func WriteXlHead(ex *sky.File, styles map[string]int, vars QuoteVars_t) error {
	if e := SetXlStyled(ex, styles, sheet, nameCell, ClientName(vars), xlStyleClient); e != nil { return e }

	if line := BirthLine(vars); line != `` {
		if e := SetXlCell(ex, sheet, birthCell, line, 0); e != nil { return e }
	}
	if line := SickCoverLine(vars); line != `` {
		if e := SetXlCell(ex, sheet, sickCell, line, 0); e != nil { return e }
	}

	if vars.core.segment != segmentEmployee {
		if e := SetXlStyled(ex, styles, sheet, payCell, `Your monthly cost`, xlStyleYouPay); e != nil { return e }
	}
	return nil
}

func DeleteXlRows(ex *sky.File, tab string, low, high int) error {
	if ex == nil { return Error(`nil excel file`) }
	if low < 1 || high < low { return nil }
	for n := low; n <= high; n++ {
		if e := ex.RemoveRow(tab, low); e != nil { return e }
	}
	return nil
}

func EnforceSlimLayout(ex *sky.File, vars QuoteVars_t, slim bool) error {
	if slim { return DeleteXlRows(ex, sheet, slimRowLow, slimRowHigh) }
	if vars.core.segment == segmentStudent && vars.core.sickCover == 0 {
		return DeleteXlRows(ex, sheet, studentNoSickRowLow, studentNoSickRowHigh)
	}
	return nil
}

func SectionItems(sec int, slim bool) []BenSecItem_t {
	var out []BenSecItem_t
	for _, item := range App.lookup.benSecItems.All() {
		if item.section != sec { continue }
		if slim && !item.isSlim { continue }
		out = append(out, item)
	}
	return out
}

func MergeXl(ex *sky.File, from, to string, style int) error {
	if ex == nil { return Error(`nil excel file`) }
	if e := ex.MergeCell(sheet, from, to); e != nil { return e }
	if style != 0 {
		if e := ex.SetCellStyle(sheet, from, to, style); e != nil { return e }
	}
	return nil
}

func WriteXlBenefits(ex *sky.File, styles map[string]int, slim bool) (lastRow int, err error) {
	row := benefitRow

	for _, sec := range App.lookup.benSecs.All() {
		items := SectionItems(sec.section, slim)
		if len(items) == 0 { continue }

		titleName := Str(sec.label, `Title`)
		labelName := Str(sec.label, `Label`)
		titleStyle := XlStyle(styles, titleName)

		aFrom := Str(`A`, row)
		aTo := Str(`A`, row+len(items)-1)
		if e := SetXlCell(ex, sheet, aFrom, sec.label, titleStyle); e != nil { return 0, e }
		if e := MergeXl(ex, aFrom, aTo, titleStyle); e != nil { return 0, e }

		for _, item := range items {
			sname := labelName
			if item.benefit == deductibleBenefit { sname = xlStyleDeductibleLabel }
			if e := SetXlStyled(ex, styles, sheet, Str(`B`, row), item.label, sname); e != nil { return 0, e }
			row++
		}
		row++
	}

	aFrom := Str(`A`, row)
	aTo := Str(`B`, row+2)
	tipStyle := XlStyle(styles, xlStyleTipTitle)
	if e := SetXlCell(ex, sheet, aFrom, `Helpful extra tips`, tipStyle); e != nil { return 0, e }
	if e := MergeXl(ex, aFrom, aTo, tipStyle); e != nil { return 0, e }
	return row + 2, nil
}

type XlPlanRow_t struct {
	plan Plan_t
	row QuoteSelectedRow_t
	choice PlanQuoteInfo_t
}

func SelectedXlPlans(vars QuoteVars_t) []XlPlanRow_t {
	state := QuoteStateFromQuoteVars(vars)
	selected := QuoteSelectedItems(state.quote)
	var out []XlPlanRow_t
	for _, item := range selected {
		row, ok := QuoteSelectedPlanRow(state, item)
		if !ok { continue }
		plan, ok := App.lookup.plans.byId[row.planId]
		if !ok { continue }
		choice, _ := vars.choices[ChoiceId_t(item.itemId)]
		out = append(out, XlPlanRow_t{
			plan: plan,
			row: QuoteSelectedRow_t{ item:item, row:row },
			choice: choice,
		})
	}
	return out
}

func OfferByFamily(benefit int, family FamilyId_t) string {
	s, ok := App.lookup.bensByFamily[BenFamily(benefit, int(family))]
	if !ok { return `` }
	return s
}

func OfferByAddons(benefit int, addons map[CategId_t]AddonId_t) string {
	if len(addons) == 0 { return `` }
	var ids []int
	for categId := range addons {
		ids = append(ids, int(categId))
	}
	sort.Ints(ids)
	for _, id := range ids {
		addon := addons[CategId_t(id)]
		if addon == 0 { continue }
		if s, ok := App.lookup.bensByAddon[BenAddon(benefit, int(addon))]; ok && Trim(s) != `` {
			return s
		}
	}
	return ``
}

func EuroWhole(amount EuroCent_t) string {
	return EuroFlatFromCent(amount).OutEuro()
}

func NoClaimRefund(plan Plan_t, amount EuroCent_t) string {
	if amount == 0 {
		note := Trim(plan.nc.note)
		if note == `` { return `No refund is available` }
		return Str(`No refund is available (`, note, `)`)
	}

	status := If(plan.nc.promise, `guaranteed`, `possible`)
	line := Str(`Refund of `, EuroWhole(amount), ` / year `, status)

	note := Trim(plan.nc.note)
	if note == `` { return line }
	return Str(line, ` (`, note, `)`)
}

func SectionValueStyle(sec BenSec_t, benefit int) string {
	if benefit == deductibleBenefit { return xlStyleDeductibleValue }
	return Str(sec.label, `Value`)
}

func PlanNameLine(plan Plan_t) string {
	return Str(plan.provName, ` / `, plan.name)
}

func RefCommission(amount EuroCent_t) int64 {
	if amount <= 0 { return 0 }
	return int64(amount) / 100
}

func PlanRefLine(x XlPlanRow_t) string {
	return Str(`Ref `, x.plan.planId, `/`, RefCommission(x.row.row.commission))
}

func WriteXlPlanHead(ex *sky.File, styles map[string]int, col string, x XlPlanRow_t) error {
	logo := Str(xlLogoDir, Lower(Trim(x.plan.provName)), `.jpg`)
	_ = ex.AddPicture(sheet, Str(col, logoRow), logo, nil)

	if top := Trim(x.plan.topNote); top != `` {
		if e := SetXlStyled(ex, styles, sheet, Str(col, topNoteRow), top, x.plan.topNoteStyle); e != nil { return e }
	}
	if e := SetXlCell(ex, sheet, Str(col, planNameRow), PlanNameLine(x.plan), 0); e != nil { return e }
	if e := SetXlStyled(ex, styles, sheet, Str(col, planRefRow), PlanRefLine(x), xlStyleCommission); e != nil { return e }
	return nil
}

func EuroCentText(amount EuroCent_t) string {
	return Str(amount.String(), ` €`)
}

func XlPreexByItem(vars QuoteVars_t) map[int]EuroCent_t {
	out := make(map[int]EuroCent_t)
	bag := UIBagVarsFromQuoteVars(vars)
	for _, x := range EditQPreexCharges(bag) {
		out[x.itemId] += x.applied
	}
	return out
}

func XlPreexByItemCateg(vars QuoteVars_t) map[string]EuroCent_t {
	out := make(map[string]EuroCent_t)
	bag := UIBagVarsFromQuoteVars(vars)
	for _, x := range EditQPreexCharges(bag) {
		key := Str(x.itemId, `:`, x.categId)
		out[key] += x.applied
	}
	return out
}

func XlPvnPreex(itemId int, row QuotePlan_t, byItemCateg map[string]EuroCent_t) EuroCent_t {
	sum := EuroCent_t(0)
	for _, addon := range row.addons {
		if addon.categId <= 0 { continue }
		if !Contains(Lower(Trim(addon.categ)), `pvn`) { continue }
		sum += byItemCateg[Str(itemId, `:`, addon.categId)]
	}
	return sum
}

func HICLine(hic, preex EuroCent_t) string {
	line := EuroCentText(hic)
	if preex > 0 { line = Str(line, ` + `, EuroCentText(preex)) }
	return line
}

func YearMaxShare(vars QuoteVars_t) EuroCent_t {
	year := CurrentDBDate().Year()
	if Valid(vars.core.buy) { year = vars.core.buy.Year() }
	yr, ok := App.lookup.years.byId[year]
	if !ok { return 0 }
	return yr.maxshare.ToEuroCent()
}

func EmployerPay(vars QuoteVars_t, total, pvn EuroCent_t) EuroCent_t {
	pay := (total - pvn) / 2
	if pay < 0 { pay = 0 }

	max := YearMaxShare(vars)
	if max > 0 && pay > max { pay = max }

	pay += (pvn / 2)
	if pay > total { pay = total }
	if pay < 0 { pay = 0 }
	return pay
}

func PlanCosts(row QuotePlan_t) (hic, pvn, sick EuroCent_t) {
	for _, addon := range row.addons {
		if !addon.priceOk { continue }
		amount := addon.base + addon.surcharge
		name := Lower(Trim(addon.categ))
		if Contains(name, `pvn`) {
			pvn += amount
			continue
		}
		if Contains(name, `sick`) {
			sick += amount
		}
	}
	hic = row.price - pvn - sick
	if hic < 0 { hic = 0 }
	return hic, pvn, sick
}

func SickAfterLine(row QuotePlan_t, vars QuoteVars_t) string {
	for _, addon := range row.addons {
		if !addon.priceOk { continue }
		if !Contains(Lower(Trim(addon.categ)), `sick`) { continue }
		daily := (int(vars.core.sickCover) / 4500) * 10
		if daily <= 0 { return `Not selected` }
		after := `29th`
		if addon.level%2 == 1 { after = `43rd` }
		return Str(daily, ` €/day as of `, after, ` day`)
	}
	return `Not selected`
}

type XlDepLine_t struct {
	label string
	price EuroCent_t
}

func XlDepLabel(dep Dependant_t, depId, age int) string {
	name := Trim(dep.name)
	if name == `` { name = Str(`Dependant `, depId) }
	if age > 0 { return Str(name, `'s monthly cost (age `, age, `)`) }
	return Str(name, `'s monthly cost`)
}

func XlDepState(vars QuoteVars_t, dep Dependant_t) State_t {
	state := QuoteStateFromQuoteVars(vars)
	if Valid(dep.birth) { state.quote[`birth`] = dep.birth.Format(`yyyymmdd`) }
	state.quote[`vision`] = If(dep.vision, `1`, ``)
	return state
}

func XlDepIsChild(age int) bool {
	return age > 0 && age < 21
}

func XlDepPreex(preex Preex_t, base EuroCent_t) EuroCent_t {
	if preex.amount.euro > 0 { return preex.amount.euro }
	if preex.amount.percent <= 0 || base <= 0 { return 0 }
	return EuroCent_t((int64(base) * int64(preex.amount.percent)) / 10000)
}

func XlDepPrice(vars QuoteVars_t, dep Dependant_t, item ChoiceId_t, plan Plan_t, row QuotePlan_t) (price EuroCent_t, age int, ok bool) {
	state := XlDepState(vars, dep)
	buyYear, yearAge, exactAge := PlanAges(state)
	if buyYear <= 0 || yearAge <= 0 { return 0, 0, false }
	age = yearAge
	if plan.exactAge { age = exactAge }
	if age <= 0 { return 0, 0, false }

	sickCover := StateInt(state, `sickCover`)
	suppressSurch := !plan.surcharge

	planBase, planSurch, planOk := CatPrice(buyYear, age, int(plan.planId), 0, sickCover, suppressSurch)
	if !planOk { return 0, age, false }

	total := planBase + planSurch
	isChild := XlDepIsChild(exactAge)
	baseByCateg := make(map[CategId_t]EuroCent_t)
	baseByCateg[0] = planBase

	for _, addon := range row.addons {
		if addon.addon == 0 { continue }
		name := Lower(Trim(addon.categ))
		if Contains(name, `sick`) { continue }
		if isChild && Contains(name, `pvn`) { continue }

		base, surch, priced := CatPrice(buyYear, age, int(addon.addon), addon.categId, sickCover, suppressSurch)
		if !priced { continue }
		total += base + surch
		if addon.categId > 0 { baseByCateg[addon.categId] += base }
	}

	if dep.vision {
		vision, has := VisionAmount(plan, planBase)
		if has { total += vision }
	}

	for _, px := range dep.preexByChoice[item] {
		total += XlDepPreex(px, baseByCateg[px.categ])
	}
	return total, age, true
}

func XlDepLines(vars QuoteVars_t, x XlPlanRow_t) []XlDepLine_t {
	var out []XlDepLine_t
	item := ChoiceId_t(x.row.item.itemId)

	for i, dep := range vars.dependants {
		price, age, ok := XlDepPrice(vars, dep, item, x.plan, x.row.row)
		if !ok { continue }

		out = append(out, XlDepLine_t{
			label: XlDepLabel(dep, i+1, age),
			price: price,
		})
	}
	return out
}

func WriteXlPlanCostRows(ex *sky.File, col string, x XlPlanRow_t, vars QuoteVars_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t) error {
	hic, pvn, sick := PlanCosts(x.row.row)
	preex := preexByItem[x.row.item.itemId]
	pvnForEmployer := pvn + XlPvnPreex(x.row.item.itemId, x.row.row, preexByItemCateg)
	total := x.row.row.price + preex
	if e := SetXlCell(ex, sheet, Str(col, hicRow), HICLine(hic, preex), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, pvnRow), EuroCentText(pvn), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, sickRow), EuroCentText(sick), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, sickAfterRow), SickAfterLine(x.row.row, vars), 0); e != nil { return e }

	youPay := total
	if vars.core.segment == segmentEmployee {
		youPay -= EmployerPay(vars, total, pvnForEmployer)
	}

	row := firstDepRow
	for _, dep := range XlDepLines(vars, x) {
		if row >= monthWithEmpRow { break }
		if e := SetXlCell(ex, sheet, Str(`A`, row), dep.label, 0); e != nil { return e }
		if e := SetXlCell(ex, sheet, Str(col, row), EuroCentText(dep.price), 0); e != nil { return e }
		youPay += dep.price
		row++
	}

	if e := SetXlCell(ex, sheet, Str(col, monthWithEmpRow), EuroCentText(total), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, monthYouPayRow), EuroCentText(youPay), 0); e != nil { return e }
	return nil
}

func BasicOffer(item BenSecItem_t, x XlPlanRow_t) (offer string, done bool) {
	if item.section != 0 { return ``, false }
	if item.secsort == 1 { return Str(EuroWhole(x.row.row.deduct), `/year`), true }
	if item.secsort == 2 { return NoClaimRefund(x.plan, x.row.row.noClaims), true }
	return `-`, true
}

func BenefitOffer(item BenSecItem_t, family FamilyId_t, addons map[CategId_t]AddonId_t) string {
	offer := OfferByFamily(item.benefit, family)
	if over := OfferByAddons(item.benefit, addons); over != `` { return over }
	return offer
}

func WriteXlPlanColumn(ex *sky.File, styles map[string]int, col string, x XlPlanRow_t, vars QuoteVars_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t, slim bool) error {
	if e := WriteXlPlanHead(ex, styles, col, x); e != nil { return e }
	if e := WriteXlPlanCostRows(ex, col, x, vars, preexByItem, preexByItemCateg); e != nil { return e }

	family := x.plan.familyId
	row := benefitRow
	for _, sec := range App.lookup.benSecs.All() {
		items := SectionItems(sec.section, slim)
		if len(items) == 0 { continue }
		for _, item := range items {
			offer, done := BasicOffer(item, x)
			if !done { offer = BenefitOffer(item, family, x.choice.addons) }
			if Trim(offer) == `` { offer = `-` }
			if e := SetXlStyled(ex, styles, sheet, Str(col, row), offer, SectionValueStyle(sec, item.benefit)); e != nil { return e }
			row++
		}
		row++
	}

	tipStyle := XlStyle(styles, xlStyleTipText)
	for k, tip := range App.lookup.familyTips[family] {
		if e := SetXlCell(ex, sheet, Str(col, row+k), tip, tipStyle); e != nil { return e }
	}
	return nil
}

func TrimXlPlanCols(ex *sky.File, count int) error {
	if ex == nil { return Error(`nil excel file`) }
	if count < 0 { count = 0 }
	if count >= len(xlPlanCols) { return nil }

	missing := len(xlPlanCols) - count
	if missing <= 0 { return nil }

	delColNum := 0
	delCount := 0
	if count == 0 {
		delColNum, _ = sky.ColumnNameToNumber(xlPlanCols[0]) // D
		delCount = (2 * missing) - 1 // D..L => 9 columns
	} else {
		firstUnused := xlPlanCols[count]
		firstUnusedNum, e := sky.ColumnNameToNumber(firstUnused)
		if e != nil || firstUnusedNum < 2 { return nil }
		delColNum = firstUnusedNum - 1 // separator before first unused plan
		delCount = 2 * missing
	}

	del, e := sky.ColumnNumberToName(delColNum)
	if e != nil || Trim(del) == `` || delCount <= 0 { return nil }

	for n := 0; n < delCount; n++ {
		if e := ex.RemoveCol(sheet, del); e != nil { return e }
	}
	return nil
}

func WriteXlPlans(ex *sky.File, styles map[string]int, vars QuoteVars_t, slim bool) error {
	plans := SelectedXlPlans(vars)
	preexByItem := XlPreexByItem(vars)
	preexByItemCateg := XlPreexByItemCateg(vars)
	for k, col := range xlPlanCols {
		if k >= len(plans) { break }
		if e := WriteXlPlanColumn(ex, styles, col, plans[k], vars, preexByItem, preexByItemCateg, slim); e != nil { return e }
	}
	return TrimXlPlanCols(ex, len(plans))
}

func XlSlimNote(vars QuoteVars_t) string {
	msg := `All plans include obligatory Long-Term Care insurance`
	if vars.core.sickCover > 0 {
		msg = Str(msg, ` and daily sick pay for an income of `, vars.core.sickCover.OutEuro())
	}
	return Str(msg, `.`)
}

func WriteXlFooterNote(ex *sky.File, styles map[string]int, vars QuoteVars_t, slim bool, row int) error {
	text := `Prices subject to increase in January each year. `
	if slim { text = Str(text, XlSlimNote(vars)) }
	return SetXlStyled(ex, styles, sheet, Str(`A`, row), text, xlStyleSlimNote)
}

func TrimXlDepRows(ex *sky.File, vars QuoteVars_t) error {
	if ex == nil { return Error(`nil excel file`) }
	count := len(vars.dependants)
	if count < 0 { count = 0 }
	if count > maxDeps { count = maxDeps }
	if count >= maxDeps { return nil }

	delRow := firstDepRow + count
	for n := maxDeps - count; n > 0; n-- {
		if e := ex.RemoveRow(sheet, delRow); e != nil { return e }
	}
	return nil
}

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

func WriteXlLayout(ex *sky.File, styles map[string]int, vars QuoteVars_t, slim bool) error {
	if e := WriteXlHead(ex, styles, vars); e != nil { return e }
	lastRow, e := WriteXlBenefits(ex, styles, slim)
	if e != nil { return e }
	if e = WriteXlFooterNote(ex, styles, vars, slim, lastRow+2); e != nil { return e }
	if e := WriteXlPlans(ex, styles, vars, slim); e != nil { return e }
	if vars.core.segment != segmentEmployee {
		if e := DeleteXlRows(ex, sheet, monthWithEmpRow, monthWithEmpRow); e != nil { return e }
	}
	if e := TrimXlDepRows(ex, vars); e != nil { return e }
	return EnforceSlimLayout(ex, vars, slim)
}
