package main

import . "pm/lib/date"

func CurrentDBDate() CalDate_t {
	rows := App.DB.Call(`quo_today_get`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var ymd int
	e := rows.Scan(&ymd)
	if e != nil { return 0 }
	return CalDate(ymd)
}
