create or replace table topnotel_plan
( plan       int not null
, lang       int not null
, note       varchar(50) not null
, softdel    bool not null default 0
, created    timestamp not null invisible default now()
, updated    timestamp not null invisible default now() on update now()
, primary key (plan, lang)
, foreign key (plan) references plans(plan)
, foreign key (lang) references languages (lang)
, check (note != '')
);

insert into topnotel_plan (plan, lang, note) select plan, 1, note from plan_topnote;
insert into topnotel_plan (plan, lang, note) select plan, 2, 'Nur für Expats.' from plan_topnote;

delimiter ###
create or replace procedure topnotel_query()
begin
    select plan, lang, note
      from topnotel_plan 
     where softdel = 0
    order by plan, lang
     ;
end
###
delimiter ;

create or replace table tipl_family
( family     int not null
, lang       int not null
, pos        int not null
, tip        varchar(150) not null
, softdel    bool not null default 0
, created    timestamp not null invisible default now()
, updated    timestamp not null invisible default now() on update now()
, primary key (family, lang, pos)
, foreign key (family) references families (family)
, foreign key (lang) references languages (lang)
, check (pos in (1,2,3))
, check (tip != '')
);

insert into tipl_family (family, pos, lang, tip)
   select family, pos, 1, tip from family_tips where trim(tip) != '';

delimiter ###
create or replace procedure tipl_family_softdel($family int, $lang int, $pos int)
begin
    update tipl_family
       set softdel = 1
     where family = $family
       and pos = $pos
       and lang = $lang;
end
###
delimiter ;

delimiter ###
create or replace procedure tipl_upsert($family int, $lang int, $pos int, $tip varchar(150))
begin
    set $tip = trim(ifnull($tip,''));

    if $tip = '' then
        call tipl_family_softdel($family, $lang, $pos);
    else
        insert into tipl_family (family, lang, pos, tip)
        values ($family, $lang, $pos, $tip)
        on duplicate key update
            tip = $tip
        ;
    end if;
end
###
delimiter ;

delimiter ###
create or replace procedure tipl_query()
begin
    select family, lang, pos, lang
      from tipl_family 
     where softdel = 0
    order by family, pos, lang
     ;
end
###
delimiter ;

insert into tipl_family (family, pos, lang, tip)
  select  31 family,   1 pos,    2 lang, "Preisgünstiger Schutz" tip union all
  select  31,   2,    2, "Hausarztprinzip" union all
  select  31,   3,    2, "Rechtsschutz bei ärztlichen Behandlungsfehlern" union all
  select  45,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select  45,   2,    2, "Große Beitragsrückerstattung" union all
  select  45,   3,    2, "Rechtsschutz bei ärztlichen Behandlungsfehlern" union all
  select  52,   1,    2, "Moderner hochwertiger Tarif"  union all
  select  52,   2,    2, "Fantastische Beitragsrückerstattung" union all
  select  52,   3,    2, "Rechtsschutz bei ärztlichen Behandlungsfehlern" union all
  select  60,   1,    2, "Nur Basisschutz" union all
  select  60,   2,    2, "Zahnbaustein optional (kann hinzugefügt oder entfernt werden)" union all
  select  60,   3,    2, "Hausarztprinzip" union all
  select  73,   1,    2, "Hausarztprinzip, auch Videosprechstunden" union all
  select  73,   2,    2, "Gehört zur AXA-Gruppe (weltweit führende Versicherungsgruppe)" union all
  select  73,   3,    2, "Zahnschutz kann ausgeschlossen werden" union all
  select  76,   1,    2, "Umfassender Basisschutz"  union all
  select  76,   2,    2, "Hausarztprinzip" union all
  select  76,   3,    2, "Großzügige spätere Tarifaufwertung"  union all
  select  87,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select  87,   2,    2, "Hausarztprinzip" union all
  select  87,   3,    2, "Beitragsbefreiung bei längerer stationärer Behandlung" union all
  select  91,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select  91,   2,    2, "Beitragsbefreiung bei längerer stationärer Behandlung" union all
  select  91,   3,    2, "Großzügige spätere Tarifaufwertung"  union all
  select  95,   1,    2, "Großzügige Rückerstattungspraxis" union all
  select  95,   2,    2, "Hochwertiger Schutz" union all
  select  95,   3,    2, "Beitragsbefreiung bei längerer stationärer Behandlung" union all
  select  99,   1,    2, "Großzügige Rückerstattungspraxis" union all
  select  99,   2,    2, "Unbegrenzter weltweiter Schutz"  union all
  select  99,   3,    2, "Rundum-Sorglos-Paket" union all
  select 104,   1,    2, "Große Beitragsrückerstattung"  union all
  select 104,   2,    2, "Hochwertiger Tarif" union all
  select 104,   3,    2, "Zweitgrößter Anbieter im deutschen Markt" union all
  select 119,   1,    2, "Solider Basisschutz"  union all
  select 119,   2,    2, "Option auf Tarifaufwertung nach 4 und 6 Jahren" union all
  select 119,   3,    2, "Zweitgrößter Anbieter im deutschen Markt" union all
  select 124,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 124,   2,    2, "Zahnleistungen gefährden die Beitragsrückerstattung nicht" union all
  select 124,   3,    2, "Flexibel anpassbar"  union all
  select 140,   1,    2, "Sehr einfacher Basisschutz" union all
  select 140,   2,    2, "Hausarztprinzip" union all
  select 140,   3,    2, "Für kostenbewusste Kunden"  union all
  select 169,   1,    2, "Umfassend, aber nur Basisschutz" union all
  select 169,   2,    2, "Hausarztprinzip" union all
  select 169,   3,    2, "10 % Rabatt für Nichtraucher mit guten Blutwerten/BMI"  union all
  select 183,   1,    2, "Spezialtarif für Expats, die vorübergehend in Deutschland sind" union all
  select 183,   2,    2, "Unkomplizierte Antragstellung - keine Gesundheitsprüfung nötig" union all
  select 183,   3,    2, "Wechsel in einen Langzeittarif bei INTER nach Gesundheitsprüfung möglich"  union all
  select 259,   1,    2, "10 % Rabatt für Nichtraucher mit guten Blutwerten/BMI"  union all
  select 259,   2,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 259,   3,    2, "Medizinische Videosprechstunde mitversichert"  union all
  select 262,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 262,   2,    2, "Hohe Beitragsrückerstattung - aber nach jeder Leistungsinanspruchnahme müssen Sie 24 Monate leistungsfrei sein, um den Bonus erneut zu erhalten."  union all
  select 262,   3,    2, "Professionelle Zahnreinigung (bis zu 120 €/Jahr) gefährdet die Beitragsrückerstattung nicht."  union all
  select 277,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 277,   2,    2, "Benutzerfreundliche App"  union all
  select 277,   3,    2, "Hallesche ist seit über 85 Jahren am Markt" union all
  select 288,   1,    2, "Benutzerfreundliche App"  union all
  select 288,   2,    2, "Digitale Gesundheitsanwendungen sind zuschussfähig," union all
  select 288,   3,    2, "Seit über 85 Jahren am Markt"  union all
  select 293,   1,    2, "Für kostenbewusste Kunden"  union all
  select 293,   2,    2, "Hausarztprinzip" union all
  select 293,   3,    2, "Benutzerfreundliche App"  union all
  select 297,   1,    2, "Benutzerfreundliche App"  union all
  select 297,   2,    2, "Hausarztprinzip" union all
  select 297,   3,    2, "Hallesche ist seit über 85 Jahren am Markt" union all
  select 301,   1,    2, "Nur für Expats mit befristeter Aufenthaltserlaubnis"  union all
  select 301,   2,    2, "Eher schwacher Schutz - für bis zu 5 Jahre"  union all
  select 301,   3,    2, "Keine Mindestvertragslaufzeit"  union all
  select 308,   1,    2, "Nur für Expats mit befristeter Aufenthaltserlaubnis"  union all
  select 308,   2,    2, "Umfassender Schutz für bis zu 5 Jahre" union all
  select 308,   3,    2, "Keine Mindestvertragslaufzeit"  union all
  select 317,   1,    2, "Nur für Expats mit befristeter Aufenthaltserlaubnis verfügbar"  union all
  select 317,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 317,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 326,   1,    2, "Hausarztprinzip" union all
  select 326,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 326,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 329,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 329,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 329,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 332,   1,    2, "Hochwertiger Schutz" union all
  select 332,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 332,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 345,   1,    2, "Leistungsaufwertung im Unfallfall: (halb-)privates Zimmer und Chefarzt"  union all
  select 345,   2,    2, "Umfassender Basistarif"  union all
  select 345,   3,    2, "Option auf Tarifaufwertung nach 3, 5 und 10 Jahren"  union all
  select 349,   1,    2, "Leistungsaufwertung im Unfallfall: privates Zimmer"  union all
  select 349,   2,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 349,   3,    2, "Option auf Tarifaufwertung nach 3, 5 und 10 Jahren"  union all
  select 353,   1,    2, "Hochwertiger Schutz" union all
  select 353,   2,    2, "Solider mittelgroßer Anbieter" union all
  select 353,   3,    2, "Halber Selbstbehalt für Kinder"  union all
  select 364,   1,    2, "Hochwertiger Schutz" union all
  select 364,   2,    2, "Verordnete digitale Gesundheits-Apps sind zuschussfähig" union all
  select 364,   3,    2, "Haushaltshilfe nach einer Operation" union all
  select 392,   1,    2, "Attraktiver Gesundheits- und Verhaltensbonus"  union all
  select 392,   2,    2, "Beitragsbefreiung bei höchstem Pflegebedarf"  union all
  select 392,   3,    2, "Außergewöhnlich niedrige Kinderbeiträge"  union all
  select 403,   1,    2, "Hausarztprinzip" union all
  select 403,   2,    2, "Der Selbstbehalt von 480 € gilt nicht für Zahnleistungen" union all
  select 403,   3,    2, "Attraktive Beiträge für Kinder"  union all
  select 419,   1,    2, "Großzügiges Tarifwechselrecht"  union all
  select 419,   2,    2, "Halber Selbstbehalt für Kinder"  union all
  select 419,   3,    2, "Telemedizin inklusive (z. B. Videoanrufe)" union all
  select 431,   1,    2, "Erfahrung seit 1843" union all
  select 431,   2,    2, "Weltweiter Schutz: bis zu 36 Monate pro Reise nach 3 Jahren" union all
  select 431,   3,    2, "Preisgünstiger Schutz" union all
  select 439,   1,    2, "Unglaublich hohe Beitragsrückerstattung" union all
  select 439,   2,    2, "Außergewöhnlich hochwertiger Schutz" union all
  select 439,   3,    2, "DKV ist der zweitgrößte Anbieter im deutschen Markt"  union all
  select 445,   1,    2, "Moderner starker Schutz" union all
  select 445,   2,    2, "Benutzerfreundliche App"  union all
  select 445,   3,    2, "Der monatliche Bonus wird jährlich mit Leistungsansprüchen verrechnet (jährlich 1.200 €)"  union all
  select 463,   1,    2, "Großzügiges Tarifwechselrecht"  union all
  select 463,   2,    2, "Solider etablierter Anbieter " union all
  select 463,   3,    2, "Universa ist seit über 175 Jahren am Markt" union all
  select 853,   1,    2, "Leistungsstarker Qualitätstarif für höchste Ansprüche"  union all
  select 853,   2,    2, "Weltweiter Schutz"  union all
  select 853,   3,    2, "Präventionszuschüsse, z. B. für Fitnessstudio-Beiträge" union all
  select 1090,   1,    2, "Hausarztprinzip" union all
  select 1090,   2,    2, "Weltweiter Schutz"  union all
  select 1090,   3,    2, "Verordnete digitale Gesundheits-Apps sind zuschussfähig" union all
  select 1091,   1,    2, "Weltweiter Schutz"  union all
  select 1091,   2,    2, "Halber Selbstbehalt für Kinder"  union all
  select 1091,   3,    2, "Verordnete digitale Gesundheits-Apps sind zuschussfähig" union all
  select 2876,   1,    2, "Guter moderner Tarif"  union all
  select 2876,   2,    2, "Benutzerfreundliche App"  union all
  select 2876,   3,    2, "Hallesche ist seit über 85 Jahren am Markt" union all
  select 2877,   1,    2, "Einfacher, aber moderner Tarif" union all
  select 2877,   2,    2, "Benutzerfreundliche App"  union all
  select 2877,   3,    2, "Hallesche ist seit über 85 Jahren am Markt" union all
  select 3640,   1,    2, "Individuelle Anpassung ist nicht möglich" union all
  select 3640,   2,    2, "Die Leistungen unterliegen staatlich vorgegebenen Budgetgrenzen" union all
  select 3640,   3,    2, "Angehörige können beitragsfrei mitversichert werden" union all
  select 3712,   1,    2, "Nur für Angestellte mit befristeter Aufenthaltserlaubnis" union all
  select 3712,   2,    2, "Wechsel in einen Langzeittarif möglich" union all
  select 3712,   3,    2, "Nur innerhalb der ersten 5 Jahre nach Einreise nach Deutschland möglich" union all
  select 3812,   1,    2, "Erfüllt höchste Ansprüche" union all
  select 3812,   2,    2, "Großzügige Optionen zur Tarifauf- und Tarifabstufung" union all
  select 3812,   3,    2, "Weltweiter Schutz bis zu 12 Monate/Reise"  union all
  select 3969,   1,    2, "Moderner hochwertiger Schutz"  union all
  select 3969,   2,    2, "Bekannte Marke"  union all
  select 3969,   3,    2, "Weltweiter Schutz für bis zu 12 Monate pro Reise"  union all
  select 3970,   1,    2, "Moderner hochwertiger Schutz für Kostenbewusste"  union all
  select 3970,   2,    2, "Bekannte Marke"  union all
  select 3970,   3,    2, "Weltweiter Schutz für bis zu 6 Monate pro Reise" union all
  select 4020,   1,    2, "Erfahrung seit 1843" union all
  select 4020,   2,    2, "Babys sind in den ersten 6 Lebensmonaten beitragsfrei" union all
  select 4020,   3,    2, "Weltweiter Schutz: bis zu 36 Monate pro Reise nach 3 Jahren" union all
  select 4269,   1,    2, "Tarif mit gutem Preis-Leistungs-Verhältnis"  union all
  select 4269,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 4269,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 4270,   1,    2, "Hochwertiger Schutz" union all
  select 4270,   2,    2, "Der einzige digitale Anbieter im deutschen Markt - kundenfreundliche App" union all
  select 4270,   3,    2, "Junges Unternehmen, das die Gewinnschwelle noch nicht erreicht hat, aber starke Investoren im Rücken hat" union all
  select 5391,   1,    2, "Moderner hochwertiger Tarif"  union all
  select 5391,   2,    2, "Zahnbaustein optional (kann hinzugefügt oder entfernt werden)" union all
  select 5391,   3,    2, "Gehört zur AXA-Gruppe (weltweit führende Versicherungsgruppe)";
 