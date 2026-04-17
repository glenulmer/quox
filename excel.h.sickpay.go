package main

import (
	. "quo2/lib/dec2"
	. "quo2/lib/output"
)

type SickPayInfo_t struct {
	daily EuroCent_t
	after string
}

func SickPayDailyFromCover(sickCover int) EuroCent_t {
	if sickCover <= 0 { return 0 }
	cover := EuroFlat_t(sickCover).ToEuroCent()
	return EuroCent_t((int64(cover) / 450000) * 1000)
}

func SickPayAfterFromLevel(level int) string {
	if level <= 0 { return `` }
	if level%2 == 1 { return `43rd` }
	return `29th`
}

func SickPayInfo(sickCover int, sickLevel int, selected bool) SickPayInfo_t {
	if !selected { return SickPayInfo_t{} }
	return SickPayInfo_t{
		daily: SickPayDailyFromCover(sickCover),
		after: SickPayAfterFromLevel(sickLevel),
	}
}

func (x SickPayInfo_t) Text() string {
	if x.daily <= 0 { return `Not selected` }
	if x.after == `` { return Str(x.daily.OutEuro(), `/day`) }
	return Str(x.daily.OutEuro(), `/day as of `, x.after, ` day`)
}

func QuotePlanSickLevel(row QuotePlan_t) (int, bool) {
	for _, addon := range row.addons {
		if addon.categId != catSick { continue }
		if addon.addon == 0 { return 0, false }
		if addon.level <= 0 { return 0, false }
		return addon.level, true
	}
	return 0, false
}

func QuotePlanSickPayInfo(row QuotePlan_t, sickCover int) SickPayInfo_t {
	level, ok := QuotePlanSickLevel(row)
	return SickPayInfo(sickCover, level, ok)
}

func QuotePlanSickPayText(row QuotePlan_t, sickCover int) string {
	return QuotePlanSickPayInfo(row, sickCover).Text()
}

func QuotePlanSickPayTextForState(state State_t, row QuotePlan_t) string {
	return QuotePlanSickPayText(row, StateInt(state, `sickCover`))
}
