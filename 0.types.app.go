package main

import "github.com/alexedwards/scs/v2"
import "net/http"
import . "pm/lib/wrapdb"
import . "pm/lib/dec2"

type YAP_t struct { year, age, productId int }
func YAP(y, a, p int) YAP_t { return YAP_t{ y, a, p } }

type Price_t struct { base, surcharge EuroCent_t }

type App_t struct {
	DB            *DB_t
	port          string
	staticVersion string
	Auth          func(http.HandlerFunc) http.HandlerFunc
	sessionManager *scs.SessionManager
	lookup struct {
		years		IdMap_t[YearVars_t]
		categs		IdMap_t[Categ_t]
		levels 		IdMap_t[Level_t]
		plans		IdMap_t[Plan_t]
		products	map[ProductId_t]Product_t
		prices		map[YAP_t]Price_t
		planAddons	map[PlanCateg_t]CatChoice_t
		planAddonChoices map[PlanCateg_t][]CatChoice_t
	}
}

var App App_t
