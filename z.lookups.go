package main

import (
	"database/sql"
	. "quo2/lib/dec2"
)

func (x YearVars_t)maxCover() EuroFlat_t { return x.cover * 2 }

func LoadStaticData() {
	App.lookup.years = LoadYearVarsIdMap()
	App.lookup.categs = LoadCategIdMap()
	App.lookup.levels = LoadLevelIdMap()
	App.lookup.products = LoadProducts()
	App.lookup.prices = LoadPrices()

	App.lookup.plans = LoadPlanDetailsIdMap()
	App.lookup.planAddons, App.lookup.planAddonChoices = LoadPlanAddons()

	App.lookup.benSecs = LoadBenSecs()
	App.lookup.benSecItems = LoadBenSecItems()
	App.lookup.bensByFamily = LoadBensByFamily()
	App.lookup.bensByAddon = LoadBensByAddon()

	App.lookup.familyTips = LoadFamilyTips()
}

type Categ_t struct {
	categId CategId_t
	name string
	catsur int
	required int
	display int
}

func LoadFamilyTips() map[FamilyId_t][]string {
	out := make(map[FamilyId_t][]string)

	rows := App.DB.Call(`quo_family_tips_query`)
	if rows.HasError() { panic(rows.Message()) }

	var family FamilyId_t
	var tip string

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&family, &tip)
		if rows.HasError() { panic(rows.Message()) }
		ftips := out[family]
		out[family] = append(ftips, tip)
	}
	if rows.HasError() { panic(rows.Message()) }

	return out
}

func LoadCategIdMap() IdMap_t[Categ_t] {
	out := IdMap[Categ_t]()
	rows := App.DB.Call(`quo_categs_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x Categ_t
		x.display = 1
		rows.Scan(&x.categId, &x.name, &x.catsur, &x.required)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(int(x.categId), x)
	}
	if rows.HasError() { panic(rows.Message()) }
	return out
}

func PlanNCCategs(planId PlanId_t) []CategId_t {
	var out []CategId_t

	rows := App.DB.Call(`quo_plan_nccategs_query`, planId)
	if rows.HasError() { panic(rows.Message()) }

	var categ CategId_t
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&categ)
		if rows.HasError() { panic(rows.Message()) }
		out = append(out, categ)
	}
	if rows.HasError() { panic(rows.Message()) }

	return out
}

func LoadPlanDetailsIdMap() IdMap_t[Plan_t] {
	out := IdMap[Plan_t]()
	rows := App.DB.Call(`quo_plan_details_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var p Plan_t
		rows.Scan(
			&p.planId, &p.familyId,
			&p.hospital, &p.dental,
			&p.priorcov, &p.noexam, &p.specref,
			&p.tempvisa, &p.surcharge, &p.shi,
			&p.vision.percent, &p.vision.euro, 
			&p.comonths,
			&p.ded.adult.euro, &p.ded.adult.percent, &p.ded.child.euro, &p.ded.child.percent,
			&p.nc.promise, &p.nc.note,
			&p.nc.adult.months, &p.nc.adult.flat, &p.nc.child.months, &p.nc.child.flat, 
			&p.name, &p.provName, &p.exactAge, &p.segmask,
			&p.topNote, &p.topNoteStyle,
		)
		if rows.HasError() { panic(rows.Message()) }
		p.ncCategs = PlanNCCategs(p.planId)
		out.Add(int(p.planId), p)
	}
	if rows.HasError() { panic(rows.Message()) }
	return out
}

type LevelId_t int
type Level_t struct {
	levelId LevelId_t
	label string
	categId int
	segments int // bitmask
	canStack bool
}

func LoadLevelIdMap() IdMap_t[Level_t] {
	out := IdMap[Level_t]()
	rows := App.DB.Call(`klec_levels_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var nostring string
	for rows.Next() {
		var x Level_t
		rows.Scan(&x.levelId, &x.label, &nostring, &x.categId, &x.segments, &x.canStack)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(int(x.levelId), x)
	}
	if rows.HasError() { panic(rows.Message()) }
	return out
}

type YearVars_t struct {
	year int
	maxshare EuroFlat_t
	cover EuroFlat_t
	ltccap EuroFlat_t
}

func LoadYearVarsIdMap() IdMap_t[YearVars_t] {
	out := IdMap[YearVars_t]()

	rows := App.DB.Call(`klec_year_get` , 0)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x YearVars_t
		var exists, isPast bool
		rows.Scan(&x.year, &x.maxshare, &x.cover, &x.ltccap, &exists, &isPast, new(sql.NullString))
		if rows.HasError() { panic(rows.Message()) }
		if isPast || !exists || x.year <= 0 { continue }
		out.Add(x.year, x)
	}
	if rows.HasError() { panic(rows.Message()) }
	return out
}

func LoadPrices() map[YAP_t]Price_t {
	prices := make(map[YAP_t]Price_t)

	rows := App.DB.Call(`quo_product_prices`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var yap YAP_t
	var pr Price_t
	for rows.Next() {
		rows.Scan(&yap.year, &yap.age, &yap.productId, &pr.base, &pr.surcharge)
		if rows.HasError() { panic(rows.Message()) }
		prices[yap] = pr
	}
	if rows.HasError() { panic(rows.Message()) }

	return prices
}

type ProductId_t int
type PlanId_t ProductId_t
type AddonId_t ProductId_t
type CategId_t int
type FamilyId_t int

func LoadProducts() map[ProductId_t]Product_t {
	products := make(map[ProductId_t]Product_t)
	rows := App.DB.Call(`quo_product_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var p Product_t
	for rows.Next() {
		rows.Scan(&p.productId, &p.providerId, &p.name, &p.categ, &p.level, &p.segmask)
		if rows.HasError() { panic(rows.Message()) }
		products[p.productId] = p
	}
	if rows.HasError() { panic(rows.Message()) }
	return products
}

type PlanCateg_t struct {
	plan	PlanId_t
	categ	CategId_t
}

type CatChoice_t struct {
	addon	AddonId_t
	level	int
	label	string
}

func LoadPlanAddons() (map[PlanCateg_t]CatChoice_t, map[PlanCateg_t][]CatChoice_t) {
	defaults := make(map[PlanCateg_t]CatChoice_t)
	choices := make(map[PlanCateg_t][]CatChoice_t)

	rows := App.DB.Call(`quo_plan_categ_addons`, 0)
	defer rows.Close()
	if rows.HasError() { panic(rows.Message()) }
	var k PlanCateg_t
	var v CatChoice_t 
	for rows.Next() {
		rows.Scan(&k.plan, &v.addon, &k.categ, &v.level, &v.label)
		if rows.HasError() { panic(rows.Message()) }
		choices[k] = append(choices[k], v)
		if _, has := defaults[k]; !has { defaults[k] = v }
	}
	if rows.HasError() { panic(rows.Message()) }

	return defaults, choices
}

type BenSec_t struct {
	section int
	label string
}

func LoadBenSecs() IdMap_t[BenSec_t] {
	out := IdMap[BenSec_t]()

	rows := App.DB.Call(`quo_bensections_query`)

	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x BenSec_t
		rows.Scan(&x.section, &x.label)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(x.section, x)
	}
	if rows.HasError() { panic(rows.Message()) }

	return out
}

type BenSecItem_t struct {
	section int
	secsort int
	benefit int
	label string
	isSlim bool
}

func LoadBenSecItems() IdMap_t[BenSecItem_t] {
	out := IdMap[BenSecItem_t]()
	seq := 0
	for secId, _ := range App.lookup.benSecs.All() {
		rows := App.DB.Call(`quo_bensecitems_query`, secId)
		if rows.HasError() { panic(rows.Message()) }

		for rows.Next() {
			var x BenSecItem_t
			x.section = secId
			rows.Scan(&x.secsort, &x.benefit, &x.label, &x.isSlim)
			if rows.HasError() { panic(rows.Message()) }
			seq++
			out.Add(seq, x)
		}
		if rows.HasError() { panic(rows.Message()) }
		rows.Close()
	}

	return out
}

type BenFamily_t struct { benefit, family int }
type BenAddon_t struct { benefit, addon int }
func BenFamily(b, f int) BenFamily_t { return BenFamily_t{ benefit:b, family:f } }
func BenAddon(b, a int) BenAddon_t { return BenAddon_t{ benefit:b, addon:a } }

func LoadBensByFamily() map[BenFamily_t]string {
	m := make(map[BenFamily_t]string)

	rows := App.DB.Call(`quo_benefits_family_query`)
	if rows.HasError() { panic(rows.Message()) }

	defer rows.Close()
	var x BenFamily_t
	for rows.Next() {
		var s string
		rows.Scan(&x.benefit, &x.family, &s)
		if rows.HasError() { panic(rows.Message()) }
		m[x] = s
	}
	if rows.HasError() { panic(rows.Message()) }

	//Log(m)
	return m
}

func LoadBensByAddon() map[BenAddon_t]string {
	m := make(map[BenAddon_t]string)

	rows := App.DB.Call(`quo_benefits_addon_query`)
	if rows.HasError() { panic(rows.Message()) }

	defer rows.Close()
	var x BenAddon_t
	for rows.Next() {
		var s string
		rows.Scan(&x.benefit, &x.addon, &s)
		if rows.HasError() { panic(rows.Message()) }
		m[x] = s
	}
	if rows.HasError() { panic(rows.Message()) }

	//Log(m)
	return m
}
