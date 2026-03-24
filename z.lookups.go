package main

import (
	"database/sql"
	. "pm/lib/dec2"
//	. "pm/lib/output"
)

func (x YearVars_t)maxCover() EuroFlat_t { return x.cover * 2 }

func LoadStaticData() {
	App.lookup.years = LoadYearVarsIdMap()
	App.lookup.categs = LoadCategIdMap()
	App.lookup.levels = LoadLevelIdMap()
	App.lookup.plans = LoadPlanDetailsIdMap()
	App.lookup.products = LoadProducts()
	App.lookup.prices = LoadPrices()
	App.lookup.planAddons = LoadPlanAddons()
}

type Categ_t struct {
	categId CategId_t
	name string
	catsur int
	required int
}

func LoadCategIdMap() IdMap_t[Categ_t] {
	out := IdMap[Categ_t]()
	rows := App.DB.Call(`quo_categs_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x Categ_t
		rows.Scan(&x.categId, &x.name, &x.catsur, &x.required)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(int(x.categId), x)
	}
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
			&p.ded.adult.euro, &p.ded.adult.percent, &p.ded.child.euro, &p.ded.adult.percent,
			&p.nc.promise, &p.nc.note,
			&p.nc.adult.months, &p.nc.adult.flat, &p.nc.child.months, &p.nc.child.flat, 
			&p.name, &p.provName, &p.exactAge, &p.segmask,
		)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(int(p.planId), p)
	}
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
	for rows.Next() {
		var x Level_t
		rows.Scan(&x.levelId, &x.label, &x.categId, &x.segments, &x.canStack)
		if rows.HasError() { panic(rows.Message()) }
		out.Add(int(x.levelId), x)
	}
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

	return prices
}

type ProductId_t int
type PlanId_t ProductId_t
type AddonId_t ProductId_t
type CategId_t int

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
	return products
}

type PlanCateg_t struct {
	plan	PlanId_t
	categ	CategId_t
}

type CatChoice_t struct {
	addon	AddonId_t
	level	int
	isdef	bool
	label	string
}

func LoadPlanAddons() map[PlanCateg_t]CatChoice_t {
	choices := make(map[PlanCateg_t]CatChoice_t)

	rows := App.DB.Call(`quo_plan_categ_addons`)
	defer rows.Close()
	if rows.HasError() { panic(rows.Message()) }
	var k PlanCateg_t
	var v CatChoice_t 
	for rows.Next() {
		rows.Scan(&k.plan, &k.categ, &v.addon, &v.level, &v.isdef, &v.label)
		if rows.HasError() { panic(rows.Message()) }
		choices[k] = v
	}

	return choices
}
