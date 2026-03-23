package main

import "net/http"
import "github.com/alexedwards/scs/v2"
import . "pm/lib/wrapdb"
import . "pm/lib/dec2"

type Sessions_t struct {
	manager *scs.SessionManager
	header  string
}

type YAP_t struct { year, age, productId int }
func YAP(y, a, p int) YAP_t { return YAP_t{ y, a, p } }

type Price_t struct { base, surcharge EuroCent_t }

type App_t struct {
	DB            *DB_t
	port          string
	staticVersion string
	Auth          func(http.HandlerFunc) http.HandlerFunc
	sessions      Sessions_t
	lookup struct {
		years		IdMap_t[YearVars_t]
		categs		IdMap_t[Categ_t]
		levels 		IdMap_t[Level_t]
		plans       []PlanName_t // alphabetical sort
		products	map[ProductId_t]Product_t
		prices		map[YAP_t]Price_t
		filters		map[PlanId_t]Filters_t
		planAddons	map[PlanCateg_t]CatChoice_t
	}
}

var App App_t
