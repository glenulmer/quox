package main

import (
	. "pm/lib/dec2"
)

type Provider_t struct {
	providerId         int
	name               string
	exact_age, display bool
}

type Family_t struct {
	familyId   int
	providerId int
	name       string
}

type SegBits_t int

type Product_t struct {
	productId  ProductId_t
	providerId int
	exactAge   bool
	name       string
	level      int
	categ      int
	segmask    SegBits_t
}

type Plan_t struct {
	exactAge bool
	familyId                      int
	hospital, dental              int
	referral, priorcov, fasttrack int
	tempvisa, register            bool
	vis                           Vision_t
	surcharge                     bool
	comonths                      Months_t // months of commission paid
	deductible                    Deductible_t
	noclaims                      NoClaims_t
	// excluding plan_noclaim_categs, since they seem to be one-off & sparse
	// excluding plan_incentives, since we have no meaningful content
}

type DeductibleSide_t struct {
	euro    EuroFlat_t
	percent Percent_t
}

type Deductible_t struct {
	adult, child  struct { euro EuroFlat_t; percent Percent_t }
}

type NoClaims_t struct {
	promise bool
	adult, child struct { months Months_t; flat EuroFlat_t }
	note string
}

type Vision_t struct {
	isPct   bool
	euro    EuroCent_t
	percent Percent_t
}
