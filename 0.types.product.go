package main

import (
	. "quo2/lib/dec2"
)

type Provider_t struct {
	providerId         int
	name               string
	exact_age, display bool
}

type Family_t struct {
	familyId   FamilyId_t
	providerId int
	name       string
}

type SegBits_t int

type Product_t struct {
	productId  ProductId_t
	providerId int
	name       string
	level      int
	categ      int
	segmask    SegBits_t
}

type Plan_t struct {
	planId PlanId_t
	familyId  FamilyId_t
	hospital, dental int
	priorcov, noexam, specref int
	tempvisa, surcharge, shi bool
	vision Vision_t
	comonths Months_t // months of commission paid
	ded Deductible_t
	nc NoClaims_t
	ncCategs []CategId_t
	// excluding plan_noclaim_categs, since they seem to be one-off & sparse
	// excluding plan_incentives, since we have no meaningful content
	name, provName  string
	exactAge bool
	segmask    SegBits_t
	topNote string
}

type Deductible_t struct {
	adult, child  struct { euro EuroCent_t; percent Percent_t }
}

type NoClaims_t struct {
	promise bool
	note string
	adult, child struct { months Months_t; flat EuroCent_t }
}

type Vision_t struct {
	percent Percent_t
	euro    EuroCent_t
}
