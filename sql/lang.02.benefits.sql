drop table if exists benl_addon;
drop table if exists benl_family;
drop table if exists benl_items;

create or replace table benl_sections
( section int not null
, lang    int not null
, label   varchar(25) not null
, primary key (section, lang)
, unique (label, lang)
, foreign key (section) references benefit_sections (section)
, foreign key (lang) references languages (lang)
, check (label != '')
);

insert into benl_sections
  select s.section, lang, s.name
    from benefit_sections s
    join languages l on l.label = 'English'
   where s.name != ''
;

create or replace table benl_items
( benefit   int not null
, lang      int not null
, label     varchar(50) not null
, primary key (benefit, lang)
, foreign key (benefit) references benefits (benefit)
, foreign key (lang) references languages (lang)
, check (label != '')
);

insert into benl_items
  select b. benefit, l.lang, b.label
    from benefits b
    join languages l on l.label = 'English'
   where b.label != ''
;

create or replace table benl_family
( benefit    int not null
, family     int not null
, lang       int not null
, label      varchar(150) not null
, softdel    bool not null default 0
, created    timestamp not null invisible default now()
, updated    timestamp not null invisible default now() on update now()
, primary key (benefit, family, lang)
, foreign key (family, benefit) references benefit_family_map (family, benefit)
, foreign key (lang) references languages (lang)
, check (softdel in (0,1))
, check (label != '')
);

insert into benl_family (benefit, family, lang, label, created, updated)
  select benefit, family, lang, descrip, created, updated
    from benefit_family_map 
    join languages l on l.label = 'English'
   where descrip != ''
;

create or replace table benl_addon 
( benefit    int not null
, addon      int not null
, lang       int not null
, label      varchar(150) not null
, softdel    bool not null default 0
, created    timestamp not null invisible default now()
, updated    timestamp not null invisible default now() on update now()
, primary key (benefit, addon, lang)
, foreign key (addon, benefit) references benefit_product_map (product, benefit)
, foreign key (lang) references languages (lang)
, check (softdel in (0,1))
, check (label != '')
);

insert into benl_addon (benefit, addon, lang, label, created, updated)
  select benefit, product, lang, descrip, created, updated
    from benefit_product_map 
    join languages l on l.label = 'English'
   where descrip != ''
;


delimiter ###
create or replace procedure benl_family_softdel($benefit int, $family int, $lang int)
begin
    update benl_family
       set softdel = 1
     where benefit = $benefit
       and family = $family
       and lang = $lang;
end
###
delimiter ;
-- call benl_family_softdel(100,200,2);


delimiter ###
create or replace procedure benl_addon_softdel($benefit int, $addon int, $lang int)
begin
    update benl_addon
       set softdel = 1
     where benefit = $benefit
       and addon = $addon
       and lang = $lang;
end
###
delimiter ;
-- call benl_addon_softdel(100,300,2);


delimiter ###
create or replace procedure benl_family_upsert($benefit int, $family int, $lang int, $label varchar(150))
begin
    set $label = trim(ifnull($label,''));

    if $label = '' then
        call benl_family_softdel($benefit, $family, $lang);
    else
        insert into benl_family (benefit, family, lang, label, softdel)
        values ($benefit, $family, $lang, $label, 0)
        on duplicate key update
            label = $label
          , softdel = 0
        ;
    end if;
end
###
delimiter ;
-- call benl_family_upsert(100,200,2,'Beispiel');


delimiter ###
create or replace procedure benl_addon_upsert($benefit int, $addon int, $lang int, $label varchar(150))
begin
    set $label = trim(ifnull($label,''));

    if $label = '' then
        call benl_addon_softdel($benefit, $addon, $lang);
    else
        insert into benl_addon (benefit, addon, lang, label, softdel)
        values ($benefit, $addon, $lang, $label, 0)
        on duplicate key update
            label = $label
          , softdel = 0
        ;
    end if;
end
###
delimiter ;
-- call benl_addon_upsert(100,300,2,'Beispiel');
