package main

import "net/http"
import . "quo2/lib/wrapdb"
import . "quo2/lib/dec2"

type YAP_t struct { year, age, productId int }
func YAP(y, a, p int) YAP_t { return YAP_t{ y, a, p } }

type Price_t struct { base, surcharge EuroCent_t }

type App_t struct {
	DB            *DB_t
	port          string
	layout        string
	staticVersion string
	Auth          func(http.HandlerFunc) http.HandlerFunc
	sessionStore *SessionStore_t
	lookup struct {
		years		IdMap_t[YearVars_t]
		categs		IdMap_t[Categ_t]
		levels 		IdMap_t[Level_t]
		products	map[ProductId_t]Product_t
		prices		map[YAP_t]Price_t

		plans		IdMap_t[Plan_t]
		planAddons	map[PlanCateg_t]CatChoice_t
		planAddonChoices map[PlanCateg_t][]CatChoice_t

		familyTips  map[FamilyId_t][]string

		benSecs		IdMap_t[BenSec_t]
		benSecItems	IdMap_t[BenSecItem_t]
		bensByFamily map[BenFamily_t]string
		bensByAddon  map[BenAddon_t]string
	}
}

var App App_t
