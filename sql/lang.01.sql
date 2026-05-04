create table if not exists languages
( lang int not null primary key
, label varchar(10) not null unique
, check (lang in (1,2))
, check (label in ('English','Deutsch'))
);

insert into languages values (1,'English') on duplicate key update label = values(label);
insert into languages values (2,'Deutsch') on duplicate key update label = values(label);

delimiter ###
create or replace procedure languages_query() 
begin
    select lang, label
      from languages x
     order by lang;
end
###
delimiter ;
-- call languages_query;


delimiter ###
create or replace function segment_id($segment varchar(50))
returns int deterministic reads sql data
begin
    declare sid int default 0;
    select segment into sid
      from segments
     where code = left($segment,1);
    return sid;
end
###
delimiter ;
-- select segment_id('emp');

delimiter ###
create or replace function addon_id($prov varchar(50), $level varchar(50), $segment varchar(50))
returns int deterministic reads sql data
begin
    declare $addon int default 0;
    select product into $addon
      from products a
     where provider = provider_id($prov)
       and level = level_id($level)
       and segmask & segment_id($segment);
    return $addon;
end
###
delimiter ;
-- select addon_id('Inter', '43A', 'e');

delimiter ###
create or replace procedure `benl_section_upsert`($section int, $lang int, $label varchar(25))
begin
    insert into benl_sections (section, lang, label)
    values ($section, $lang, $label)
    on duplicate key update
        label = $label;
end
###
delimiter ;


delimiter ###
create or replace procedure `benl_item_upsert`($benefit int, $lang int, $label varchar(50))
begin
    insert into benl_items (benefit, lang, label)
    values ($benefit, $lang, $label)
    on duplicate key update
        label = $label;
end
###
delimiter ;
