package main 

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/date"
)

func di(name string, date CalDate_t) Elem_t {
	return DateInput().Name(name).Req().Value(date.Format(`yyyy-mm-dd`))
}

func ti() Elem_t { return Elem(`input`).KV(`autocomplete`,`off`).KV(`type`,`text`) }


func CustomerCard() Elem_t {
	today := CurrentDBDate()
	birth := DateFromYMD(today.Year()-32, 6, 15)
	buy := today.Days(40).ToWorkDay()
	body := Div().Class(`card-body`).
		Id(`Customer`).
		Wrap(Card(
			Field(`Name`, 12).Wrap(ti().Name(`name`)),
			Field(`Birthdate`, 6).Wrap(di(`birth`, birth)),
			Field(`Buy date`, 6).Wrap(di(`buy`, buy)),
		))
	title := Div(`Customer`).Class(`card-title`)
	return Div().Class(`card`).Wrap(title, body)
}
