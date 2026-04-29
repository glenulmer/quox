package main

import (
	"sort"

	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

const sheet = `Sheet1`
const quoteSheet = `Quote`
const summarySheet = `Summary`
const xlSummaryPlanBlue = `#E6F2FC`
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

func BirthLine(vars QuoteVars_t, lang LangId_t) string {
	return BirthDateLine(lang, vars.core.birth)
}

func SickCoverLine(vars QuoteVars_t, lang LangId_t) string {
	return SickPayIncomeLine(lang, vars.core.sickCover)
}

func WriteXlHead(ex *sky.File, styles map[string]int, vars QuoteVars_t, lang LangId_t) error {
	if e := SetXlStyled(ex, styles, sheet, nameCell, ClientName(vars), xlStyleClient); e != nil { return e }

	if line := BirthLine(vars, lang); line != `` {
		if e := SetXlCell(ex, sheet, birthCell, line, 0); e != nil { return e }
	}
	if line := SickCoverLine(vars, lang); line != `` {
		if e := SetXlCell(ex, sheet, sickCell, line, 0); e != nil { return e }
	}
	if e := SetXlCell(ex, sheet, `A6`, HealthInsuranceLine(lang), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, `A7`, LongTermCareLine(lang), 0); e != nil { return e }
	if vars.core.segment == segmentEmployee {
		if e := SetXlCell(ex, sheet, `A21`, TotalWithEmployerLine(lang), 0); e != nil { return e }
	}
	if e := SetXlCell(ex, sheet, payCell, YourMonthlyCostLine(lang), 0); e != nil { return e }

	if vars.core.segment != segmentEmployee {
		if e := SetXlStyled(ex, styles, sheet, payCell, YourMonthlyCostLine(lang), xlStyleYouPay); e != nil { return e }
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

func SummaryEuroText(amount EuroCent_t) string {
	if amount == 0 { return `-` }
	return EuroCentText(amount)
}

func SummaryDateText(x CalDate_t) string {
	if !Valid(x) { return `-` }
	return x.Format(`yyyy-mm-dd`)
}

func SummaryBoolText(x bool) string {
	return If(x, `Yes`, `No`)
}

func SummaryChoiceText(sp string, id int, args ...any) string {
	if id <= 0 { return `-` }
	for _, x := range QuoteChoices(sp, args...) {
		if x.id == id { return x.label }
	}
	return Str(id)
}

func SummaryAddonText(x QuotePlanAddon_t) string {
	if x.addon == 0 { return `` }
	return AddonName(CatChoice_t{ addon:x.addon, label:x.label })
}

func SummaryVisionPaid(row QuotePlan_t) EuroCent_t {
	for _, addon := range row.addons {
		if addon.categId > 0 { continue }
		if !Contains(Lower(Trim(addon.categ)), `vision`) { continue }
		if !addon.priceOk { return 0 }
		return addon.base + addon.surcharge
	}
	return 0
}

func XlPreexFromBase(preex Preex_t, base EuroCent_t) EuroCent_t {
	if preex.amount.euro > 0 { return preex.amount.euro }
	if preex.amount.percent <= 0 || base <= 0 { return 0 }
	return EuroCent_t((int64(base) * int64(preex.amount.percent)) / 10000)
}

func XlPlanBaseByCateg(row QuotePlan_t) map[CategId_t]EuroCent_t {
	out := make(map[CategId_t]EuroCent_t)
	out[0] = row.planBase
	for _, addon := range row.addons {
		if !addon.priceOk { continue }
		if addon.categId <= 0 { continue }
		if addon.base <= 0 { continue }
		out[addon.categId] += addon.base
	}
	return out
}

func SummaryPreexMaps(vars QuoteVars_t, plans []XlPlanRow_t) (map[int]EuroCent_t, map[string]EuroCent_t, map[string]string) {
	byItem := make(map[int]EuroCent_t)
	byItemCateg := make(map[string]EuroCent_t)
	notes := make(map[string]string)

	for _, x := range plans {
		itemId := x.row.item.itemId
		baseByCateg := XlPlanBaseByCateg(x.row.row)
		for _, preex := range x.choice.preex {
			applied := XlPreexFromBase(preex, baseByCateg[preex.categ])
			if applied <= 0 { continue }
			key := Str(itemId, `:`, preex.categ)
			byItem[itemId] += applied
			byItemCateg[key] += applied
			if Trim(preex.note) != `` && Trim(notes[key]) == `` {
				notes[key] = preex.note
			}
		}
	}
	return byItem, byItemCateg, notes
}

type XlDepCharge_t struct {
	plan string
	level string
	applied EuroCent_t
	note string
}

func XlPlanLevelLabel(row QuotePlan_t, categId CategId_t) string {
	if categId == 0 { return `Plan` }
	for _, addon := range row.addons {
		if addon.categId != categId { continue }
		level := QuoteAddonPickText(addon)
		if level == `` { return addon.categ }
		return level
	}
	return Str(categId)
}

func XlDepCharges(dep Dependant_t, plans []XlPlanRow_t) []XlDepCharge_t {
	var out []XlDepCharge_t
	for _, plan := range plans {
		item := ChoiceId_t(plan.row.item.itemId)
		preexList := dep.preexByChoice[item]
		if len(preexList) == 0 { continue }

		baseByCateg := XlPlanBaseByCateg(plan.row.row)
		for _, preex := range preexList {
			applied := XlPreexFromBase(preex, baseByCateg[preex.categ])
			if applied <= 0 { continue }
			out = append(out, XlDepCharge_t{
				plan: plan.row.row.label,
				level: XlPlanLevelLabel(plan.row.row, preex.categ),
				applied: applied,
				note: preex.note,
			})
		}
	}
	return out
}

func WriteSummaryKV(ex *sky.File, tab string, row int, key, val string, keyStyle int) error {
	if e := SetXlCell(ex, tab, Str(`A`, row), key, keyStyle); e != nil { return e }
	return SetXlCell(ex, tab, Str(`B`, row), val, 0)
}

func SummaryHeadStyle(ex *sky.File) int {
	if ex == nil { return 0 }
	style, e := ex.NewStyle(&sky.Style{
		Font: &sky.Font{ Bold:true },
	})
	if e != nil { return 0 }
	return style
}

func SummaryTitleStyle(ex *sky.File) int {
	if ex == nil { return 0 }
	style, e := ex.NewStyle(&sky.Style{
		Font: &sky.Font{ Bold:true, Size:14 },
	})
	if e != nil { return 0 }
	return style
}

func SummaryPlanStyle(ex *sky.File) int {
	if ex == nil { return 0 }
	style, e := ex.NewStyle(&sky.Style{
		Fill: sky.Fill{
			Type: `pattern`,
			Color: []string{xlSummaryPlanBlue},
			Pattern: 1,
		},
	})
	if e != nil { return 0 }
	return style
}

func SummaryPlanHeadStyle(ex *sky.File) int {
	if ex == nil { return 0 }
	style, e := ex.NewStyle(&sky.Style{
		Font: &sky.Font{ Bold:true },
		Fill: sky.Fill{
			Type: `pattern`,
			Color: []string{xlSummaryPlanBlue},
			Pattern: 1,
		},
	})
	if e != nil { return 0 }
	return style
}

func WriteXlSummary(ex *sky.File, vars QuoteVars_t) error {
	if ex == nil { return Error(`nil excel file`) }
	if ix, e := ex.GetSheetIndex(summarySheet); e == nil && ix >= 0 {
		if e := ex.DeleteSheet(summarySheet); e != nil { return e }
	}
	if _, e := ex.NewSheet(summarySheet); e != nil { return e }
	_ = ex.SetColWidth(summarySheet, `A`, `A`, 34)
	_ = ex.SetColWidth(summarySheet, `B`, `B`, 40)
	_ = ex.SetColWidth(summarySheet, `C`, `C`, 16)
	_ = ex.SetColWidth(summarySheet, `D`, `D`, 42)
	showGrid := false
	if e := ex.SetSheetView(summarySheet, 0, &sky.ViewOptions{ ShowGridLines:&showGrid }); e != nil { return e }

	titleStyle := SummaryTitleStyle(ex)
	headStyle := SummaryHeadStyle(ex)
	planStyle := SummaryPlanStyle(ex)
	planHeadStyle := SummaryPlanHeadStyle(ex)

	row := 1
	if e := SetXlCell(ex, summarySheet, `A1`, `Quote Summary`, titleStyle); e != nil { return e }
	if e := ex.SetCellStyle(summarySheet, `A1`, `D1`, titleStyle); e != nil { return e }
	row += 2

	core := []struct{ key, val string }{
		{ `Client name`, ClientName(vars) },
		{ `Segment`, SummaryChoiceText(`klpm_segments_chooser`, vars.core.segment) },
		{ `Birth date`, SummaryDateText(vars.core.birth) },
		{ `Buy date`, SummaryDateText(vars.core.buy) },
		{ `Sick cover`, vars.core.sickCover.OutEuro() },
		{ `Prior cover`, SummaryChoiceText(`klpm_priorcov_chooser`, vars.core.priorCov) },
		{ `Exam`, SummaryChoiceText(`klpm_noexam_chooser`, vars.core.exam) },
		{ `Specialist`, SummaryChoiceText(`klpm_specialist_chooser`, vars.core.specref) },
		{ `Vision`, SummaryBoolText(vars.core.vision) },
		{ `Temp visa`, SummaryBoolText(vars.core.tempVisa) },
		{ `No PVN`, SummaryBoolText(vars.core.noPVN) },
		{ `Natural medicine`, SummaryBoolText(vars.core.naturalMed) },
	}
	for _, x := range core {
		if e := WriteSummaryKV(ex, summarySheet, row, x.key, x.val, headStyle); e != nil { return e }
		row++
	}
	row++

	plans := SelectedXlPlans(vars)
	preexByItem, preexByItemCateg, preexNoteByItemCateg := SummaryPreexMaps(vars, plans)
	if len(plans) == 0 {
		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `No selected plans`, headStyle); e != nil { return e }
		row += 2
	}
	for k, plan := range plans {
		planTop := row
		preexRow := 0
		visionRow := 0
		categHeadRow := 0
		itemId := plan.row.item.itemId
		if e := SetXlCell(ex, summarySheet, Str(`A`, row), Str(`Plan `, k+1), headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`B`, row), plan.row.row.label, headStyle); e != nil { return e }
		row++
		if plan.plan.exactAge {
			if e := WriteSummaryKV(ex, summarySheet, row, `Age mode`, `Exact-age`, headStyle); e != nil { return e }
			row++
		}

		preexRow = row
		if e := WriteSummaryKV(ex, summarySheet, row, `Pre-ex total`, SummaryEuroText(preexByItem[itemId]), headStyle); e != nil { return e }
		row++
		if vars.core.vision {
			visionRow = row
			if e := WriteSummaryKV(ex, summarySheet, row, `Vision correction`, SummaryEuroText(SummaryVisionPaid(plan.row.row)), headStyle); e != nil { return e }
			row++
			row++
		}

		categHeadRow = row
		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `Category`, headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`B`, row), `Selected addon`, headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`C`, row), `Pre-ex`, headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`D`, row), `Optional note`, headStyle); e != nil { return e }
		row++

		planKey := Str(itemId, `:`, 0)
		planPreex := preexByItemCateg[planKey]
		planNote := ``
		if planPreex > 0 { planNote = preexNoteByItemCateg[planKey] }
		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `Plan`, 0); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`B`, row), ``, 0); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`C`, row), SummaryEuroText(planPreex), 0); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`D`, row), planNote, 0); e != nil { return e }
		row++

		for _, addon := range plan.row.row.addons {
			if addon.categId <= 0 { continue }
			key := Str(itemId, `:`, addon.categId)
			preex := preexByItemCateg[key]
			note := ``
			if preex > 0 { note = preexNoteByItemCateg[key] }
			if e := SetXlCell(ex, summarySheet, Str(`A`, row), addon.categ, 0); e != nil { return e }
			if e := SetXlCell(ex, summarySheet, Str(`B`, row), SummaryAddonText(addon), 0); e != nil { return e }
			if e := SetXlCell(ex, summarySheet, Str(`C`, row), SummaryEuroText(preex), 0); e != nil { return e }
			if e := SetXlCell(ex, summarySheet, Str(`D`, row), note, 0); e != nil { return e }
			row++
		}

		planLow := row - 1
		if planStyle != 0 {
			if e := ex.SetCellStyle(summarySheet, Str(`A`, planTop), Str(`D`, planLow), planStyle); e != nil { return e }
		}
		if planHeadStyle != 0 {
			if e := ex.SetCellStyle(summarySheet, Str(`A`, planTop), Str(`B`, planTop), planHeadStyle); e != nil { return e }
			if e := ex.SetCellStyle(summarySheet, Str(`A`, preexRow), Str(`A`, preexRow), planHeadStyle); e != nil { return e }
			if visionRow > 0 {
				if e := ex.SetCellStyle(summarySheet, Str(`A`, visionRow), Str(`A`, visionRow), planHeadStyle); e != nil { return e }
			}
			if e := ex.SetCellStyle(summarySheet, Str(`A`, categHeadRow), Str(`D`, categHeadRow), planHeadStyle); e != nil { return e }
		}
		row++
	}

	deps := vars.dependants
	if len(deps) > 0 {
		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `Dependants`, headStyle); e != nil { return e }
		row += 2
	}
	for _, dep := range deps {
		charges := XlDepCharges(dep, plans)
		if Trim(dep.name) == `` && !Valid(dep.birth) && !dep.vision && len(charges) == 0 { continue }

		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `Dependant`, headStyle); e != nil { return e }
		row++
		if e := WriteSummaryKV(ex, summarySheet, row, `Name`, dep.name, headStyle); e != nil { return e }
		row++
		if e := WriteSummaryKV(ex, summarySheet, row, `Birth date`, SummaryDateText(dep.birth), headStyle); e != nil { return e }
		row++
		if e := WriteSummaryKV(ex, summarySheet, row, `Vision`, SummaryBoolText(dep.vision), headStyle); e != nil { return e }
		row++

		if len(charges) == 0 {
			if e := WriteSummaryKV(ex, summarySheet, row, `Pre-ex charges`, `-`, headStyle); e != nil { return e }
			row += 2
			continue
		}

		if e := SetXlCell(ex, summarySheet, Str(`A`, row), `Plan / Category`, headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`B`, row), `Pre-ex`, headStyle); e != nil { return e }
		if e := SetXlCell(ex, summarySheet, Str(`C`, row), `Optional note`, headStyle); e != nil { return e }
		row++
		for _, x := range charges {
			if e := SetXlCell(ex, summarySheet, Str(`A`, row), Str(x.plan, ` / `, x.level), 0); e != nil { return e }
			if e := SetXlCell(ex, summarySheet, Str(`B`, row), EuroCentText(x.applied), 0); e != nil { return e }
			if e := SetXlCell(ex, summarySheet, Str(`C`, row), x.note, 0); e != nil { return e }
			row++
		}
		row++
	}

	return nil
}

func FinalizeXlSheets(ex *sky.File, vars QuoteVars_t) error {
	if ex == nil { return Error(`nil excel file`) }
	if e := ex.SetSheetName(sheet, quoteSheet); e != nil { return e }
	return WriteXlSummary(ex, vars)
}

func XlPreexByItem(vars QuoteVars_t) map[int]EuroCent_t {
	plans := SelectedXlPlans(vars)
	byItem, _, _ := SummaryPreexMaps(vars, plans)
	return byItem
}

func XlPreexByItemCateg(vars QuoteVars_t) map[string]EuroCent_t {
	plans := SelectedXlPlans(vars)
	_, byItemCateg, _ := SummaryPreexMaps(vars, plans)
	return byItemCateg
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

func SickAfterLine(row QuotePlan_t, vars QuoteVars_t, lang LangId_t) string {
	for _, addon := range row.addons {
		if !addon.priceOk { continue }
		if !Contains(Lower(Trim(addon.categ)), `sick`) { continue }
		daily := (int(vars.core.sickCover) / 4500) * 10
		if daily <= 0 { return SickPayWaitingLine(lang, 0, 0, false) }
		waitingDays := 29
		if addon.level%2 == 1 { waitingDays = 43 }
		return SickPayWaitingLine(lang, daily, waitingDays, true)
	}
	return SickPayWaitingLine(lang, 0, 0, false)
}

type XlDepLine_t struct {
	label string
	price EuroCent_t
}

func XlDepLabel(dep Dependant_t, depId, age int, lang LangId_t) string {
	return DependantMonthlyCostLine(lang, dep.name, age)
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

func XlDepLines(vars QuoteVars_t, x XlPlanRow_t, lang LangId_t) []XlDepLine_t {
	var out []XlDepLine_t
	item := ChoiceId_t(x.row.item.itemId)

	for i, dep := range vars.dependants {
		price, age, ok := XlDepPrice(vars, dep, item, x.plan, x.row.row)
		if !ok { continue }

		out = append(out, XlDepLine_t{
			label: XlDepLabel(dep, i+1, age, lang),
			price: price,
		})
	}
	return out
}

func WriteXlPlanCostRows(ex *sky.File, col string, x XlPlanRow_t, vars QuoteVars_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t, lang LangId_t) error {
	hic, pvn, sick := PlanCosts(x.row.row)
	preex := preexByItem[x.row.item.itemId]
	pvnForEmployer := pvn + XlPvnPreex(x.row.item.itemId, x.row.row, preexByItemCateg)
	total := x.row.row.price + preex
	if e := SetXlCell(ex, sheet, Str(col, hicRow), HICLine(hic, preex), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, pvnRow), EuroCentText(pvn), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, sickRow), EuroCentText(sick), 0); e != nil { return e }
	if e := SetXlCell(ex, sheet, Str(col, sickAfterRow), SickAfterLine(x.row.row, vars, lang), 0); e != nil { return e }

	youPay := total
	if vars.core.segment == segmentEmployee {
		youPay -= EmployerPay(vars, total, pvnForEmployer)
	}

	row := firstDepRow
	for _, dep := range XlDepLines(vars, x, lang) {
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

func WriteXlPlanColumn(ex *sky.File, styles map[string]int, col string, x XlPlanRow_t, vars QuoteVars_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t, slim bool, lang LangId_t) error {
	if e := WriteXlPlanHead(ex, styles, col, x); e != nil { return e }
	if e := WriteXlPlanCostRows(ex, col, x, vars, preexByItem, preexByItemCateg, lang); e != nil { return e }

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

func WriteXlPlans(ex *sky.File, styles map[string]int, vars QuoteVars_t, slim bool, lang LangId_t) error {
	plans := SelectedXlPlans(vars)
	preexByItem := XlPreexByItem(vars)
	preexByItemCateg := XlPreexByItemCateg(vars)
	for k, col := range xlPlanCols {
		if k >= len(plans) { break }
		if e := WriteXlPlanColumn(ex, styles, col, plans[k], vars, preexByItem, preexByItemCateg, slim, lang); e != nil { return e }
	}
	return TrimXlPlanCols(ex, len(plans))
}

func XlSlimNote(vars QuoteVars_t, lang LangId_t) string {
	msg := Str(`All plans include `, LongTermCareLine(lang))
	if vars.core.sickCover > 0 {
		msg = Str(msg, ` and `, SickPayIncomeLine(lang, vars.core.sickCover))
	}
	return Str(msg, `.`)
}

func WriteXlFooterNote(ex *sky.File, styles map[string]int, vars QuoteVars_t, slim bool, row int, lang LangId_t) error {
	text := `Prices subject to increase in January each year. `
	if slim { text = Str(text, XlSlimNote(vars, lang)) }
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

func WriteXlLayout(ex *sky.File, styles map[string]int, vars QuoteVars_t) error {
	lang := vars.lang
	if lang <= 0 { lang = English }
	slim := vars.slim == 1
	if e := WriteXlHead(ex, styles, vars, lang); e != nil { return e }
	lastRow, e := WriteXlBenefits(ex, styles, slim)
	if e != nil { return e }
	if e = WriteXlFooterNote(ex, styles, vars, slim, lastRow+2, lang); e != nil { return e }
	if e := WriteXlPlans(ex, styles, vars, slim, lang); e != nil { return e }
	if vars.core.segment != segmentEmployee {
		if e := DeleteXlRows(ex, sheet, monthWithEmpRow, monthWithEmpRow); e != nil { return e }
	}
	if e := TrimXlDepRows(ex, vars); e != nil { return e }
	if e := EnforceSlimLayout(ex, vars, slim); e != nil { return e }
	return FinalizeXlSheets(ex, vars)
}
