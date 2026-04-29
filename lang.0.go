package main

type LangId_t int
type Lang_t struct {
	langId LangId_t
	label string
}

const English, German LangId_t = 1, 2
