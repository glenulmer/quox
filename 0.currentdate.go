package main

import . "pm/lib/date"

const spCurrentDateQuery = `quo_current_date_query`

func CurrentDBDate() CalDate_t {
	rows := App.DB.Call(spCurrentDateQuery)
	if rows.HasError() { return 0 }
	defer rows.Close()
	if !rows.Next() { return 0 }
	var ymd int
	e := rows.Scan(&ymd)
	if e != nil { return 0 }
	out := CalDate(ymd)
	if !Valid(out) { return 0 }
	return out
}
