package main

import "strings"
import . "pm/lib/output"

const spPriorCovQuery = `klec_priorcov_query`
const spReferralsQuery = `klec_referrals_query`
const spLevelChooser = `quo_level_chooser`
const spPlanDeductiblesDistinct = `plan_deductibles_distinct`
const specialistAnyCode = 2
const examNoExamCode = 1

func CategIDByName(idMap IdMap_t[Categ_t], name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, id := range idMap.sort {
		x, ok := idMap.byId[id]
		if !ok { continue }
		if strings.ToLower(strings.TrimSpace(x.name)) == name { return x.categId }
	}
	return 0
}

func QueryPriorCoverOptions() (list []SelectOption_t) {
	rows := App.DB.Call(spPriorCovQuery)
	if rows.HasError() { panic(Error(`call `, spPriorCovQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x SelectOption_t
		e := rows.Scan(&x.id, &x.name)
		if e != nil { panic(Error(`scan `, spPriorCovQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spPriorCovQuery, ` failed: `, e)) }
	return
}

func QueryReferralOptions() (list []SelectOption_t) {
	rows := App.DB.Call(spReferralsQuery)
	if rows.HasError() { panic(Error(`call `, spReferralsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x SelectOption_t
		e := rows.Scan(&x.id, &x.name)
		if e != nil { panic(Error(`scan `, spReferralsQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spReferralsQuery, ` failed: `, e)) }
	return
}

func QueryLevelChooser(categ int) (levels []SelectOption_t) {
	rows := App.DB.Call(spLevelChooser, categ)
	if rows.HasError() { panic(Error(`call `, spLevelChooser, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x SelectOption_t
		e := rows.Scan(&x.id, &x.name)
		if e != nil { panic(Error(`scan `, spLevelChooser, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.id <= 0 || x.name == `` { continue }
		levels = append(levels, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spLevelChooser, ` failed: `, e)) }
	return
}
