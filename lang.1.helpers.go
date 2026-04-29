package main

import (
	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

func BirthDateLine(lang LangId_t, birth CalDate_t) string {
	if !Valid(birth) { return `` }
	switch lang {
	case German:
		return Str(`Geburtsdatum: `, birth.Format(`d. mon, yyyy`))
	default:
		return Str(`Date of birth: `, birth.Format(`d. mon, yyyy`))
	}
}

func HealthInsuranceLine(lang LangId_t) string {
	switch lang {
	case German:
		return `Krankenversicherung`
	default:
		return `Health insurance`
	}
}

func LongTermCareLine(lang LangId_t) string {
	switch lang {
	case German:
		return `Gesetzliche Pflegepflichtversicherung`
	default:
		return `Obligatory Long-Term Care insurance`
	}
}

func SickPayIncomeLine(lang LangId_t, sickCover EuroFlat_t) string {
	if sickCover <= 0 { return `` }
	switch lang {
	case German:
		return Str(`Krankentagegeld (bei einem Einkommen von `, sickCover.OutEuro(), `)`)
	default:
		return Str(`Daily Sick Pay (for a `, sickCover.OutEuro(), ` income)`)
	}
}

func SickPayWaitingLine(lang LangId_t, dailyEuro int, waitingDays int, selected bool) string {
	if !selected {
		switch lang {
		case German:
			return `Nicht ausgewählt`
		default:
			return `Not selected`
		}
	}

	switch lang {
	case German:
		return Str(dailyEuro, ` €/Tag ab dem `, waitingDays, `. Tag`)
	default:
		return Str(dailyEuro, ` €/day as of day `, waitingDays)
	}
}

func DependantMonthlyCostLine(lang LangId_t, name string, age int) string {
	n := Trim(name)
	if n == `` {
		switch lang {
		case German:
			n = `Mitversicherte Person`
		default:
			n = `Dependant`
		}
	}

	if age > 0 {
		switch lang {
		case German:
			return Str(`Monatlicher Beitrag von `, n, ` (Alter `, age, `)`)
		default:
			return Str(n, `'s monthly cost (age `, age, `)`)
		}
	}

	switch lang {
	case German:
		return Str(`Monatlicher Beitrag von `, n)
	default:
		return Str(n, `'s monthly cost`)
	}
}

func TotalWithEmployerLine(lang LangId_t) string {
	switch lang {
	case German:
		return `Gesamter Monatsbeitrag (inkl. Arbeitgeberzuschuss)`
	default:
		return `Total monthly cost (incl. employer subsidy)`
	}
}

func YourMonthlyCostLine(lang LangId_t) string {
	switch lang {
	case German:
		return `Ihr Monatsbeitrag (nach Arbeitgeberzuschuss)`
	default:
		return `Your monthly cost (after employer subsidy)`
	}
}
