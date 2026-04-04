package htmlHelper

import (
	. "quo2/lib/output"
)

type tHeadItem struct { kind, text string }
type Head_t []tHeadItem

func Head() Head_t { return Head_t{} }

func (x Head_t)Title(t string) Head_t   { return append(x, tHeadItem{`title`, t}) }
func (x Head_t)Icon(t string) Head_t    { return append(x, tHeadItem{`icon`, t}) }
func (x Head_t)CSS(t string) Head_t     { return append(x, tHeadItem{`css`, t}) }
func (x Head_t)CSSTail(t string) Head_t { return append(x, tHeadItem{`cTail`, t}) }
func (x Head_t)JS(t string) Head_t      { return append(x, tHeadItem{`js`, t}) }
func (x Head_t)JSTail(t string) Head_t  { return append(x, tHeadItem{`jsTail`, t}) }
func (x Head_t)Script(t string) Head_t  { return append(x, tHeadItem{`script`, t}) }
func (x Head_t)End() Head_t { return x }

func (x Head_t)Left() string {
	const tab = "\t"
	var b, ctail Builder
	b.Add(
		`<!DOCTYPE html>`, NL,
		`<html lang="en">`, NL,
		`<head>`, NL,
		`	<meta charset="UTF-8">`, NL,
		`	<meta name="viewport" content="width=device-width, initial-scale=1.0">`, NL,
		)

	for _, item := range x {
		switch item.kind {
		case "title":  b.Add(tab, `<title>`, item.text, `</title>`, NL)
		case "icon":   b.Add(tab, `<link rel="icon" href="`, item.text, `" type="image/x-icon">`, NL)
		case "css":    b.Add(tab, `<link rel="stylesheet" href="`, item.text, `">`, NL)
		case "js":     b.Add(tab, `<script src="`, item.text, `"></script>`, NL)
		case "script": b.Add(tab, `<script>`, NL, item.text, NL, tab, "</script>", NL)
		case "style":  b.Add(tab, `<style>`, NL, item.text, NL, tab, `</style>`, NL)
		case "meta":   b.Add(tab, `<meta `, item.text,`>`, NL)
		case "cTail":  ctail.Add(tab, `<link rel="stylesheet" href="`, item.text, `">`, NL)
		}
	}
	
	b.Add(ctail.String())
	b.Add(
		`</head>`, NL,
		`<body class="content">`, NL, // content => use bulma
	)

	return b.String()
}

func (x Head_t)Right() string {
	const tab = "\t"
	var b Builder
	for _, item := range x {
		switch item.kind {
		case "jsTail": b.Add(tab, `<script src="`, item.text, `"></script>`, NL)
		}
	}
	b.Add(`</body>`)
	return b.String()
}
