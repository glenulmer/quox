package main

import (
	. "quo2/lib/dec2"
	. "quo2/lib/output"
)

const segmentEmployee = 1

func QuotePVNTotal(row QuotePlan_t) EuroCent_t {
	var out EuroCent_t
	for _, addon := range row.addons {
		if !addon.priceOk { continue }
		if !Contains(Lower(Trim(addon.categ)), `pvn`) { continue }
		out += addon.base + addon.surcharge
	}
	return out
}

func EmployerShareCapByYear(buyYear int) EuroCent_t {
	if buyYear <= 0 { return 0 }
	x, ok := App.lookup.years.byId[buyYear]
	if !ok { return 0 }
	return x.maxshare.ToEuroCent()
}

func EmployerShareAmount(segment int, total, pvn, cap EuroCent_t) EuroCent_t {
	if segment != segmentEmployee { return 0 }
	if total <= 0 { return 0 }

	if pvn < 0 {
		pvn = 0
	} else if pvn > total {
		pvn = total
	}

	withoutPVN := total - pvn
	share := withoutPVN / 2
	if cap > 0 && share > cap { share = cap }
	share += pvn / 2

	if share < 0 { return 0 }
	if share > total { return total }
	return share
}

func QuoteEmployerShare(row QuotePlan_t, segment, buyYear int) EuroCent_t {
	return EmployerShareAmount(segment, row.price, QuotePVNTotal(row), EmployerShareCapByYear(buyYear))
}

func QuoteYouPay(row QuotePlan_t, segment, buyYear int) EuroCent_t {
	share := QuoteEmployerShare(row, segment, buyYear)
	you := row.price - share
	if you < 0 { return 0 }
	return you
}

func QuoteYouPayWithDependants(row QuotePlan_t, depTotal EuroCent_t, segment, buyYear int) EuroCent_t {
	you := QuoteYouPay(row, segment, buyYear) + depTotal
	if you < 0 { return 0 }
	return you
}

func QuoteEmployerShareForState(state State_t, row QuotePlan_t) EuroCent_t {
	buyYear, _, _ := PlanAges(state)
	segment := StateInt(state, `segment`)
	return QuoteEmployerShare(row, segment, buyYear)
}

func QuoteYouPayForState(state State_t, row QuotePlan_t) EuroCent_t {
	buyYear, _, _ := PlanAges(state)
	segment := StateInt(state, `segment`)
	return QuoteYouPay(row, segment, buyYear)
}
