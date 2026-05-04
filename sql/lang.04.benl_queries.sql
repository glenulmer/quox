delimiter ###
create or replace procedure klpm_bensections_query()
begin
    select section, lang, label from benl_sections b order by section, lang;
end
###
delimiter ;


delimiter ###
create or replace procedure klpm_bensecitems_query($sec int)
begin
    select b.secsort, b.benefit, b.slim, i.lang, i.label
      from benefits b
      join benl_items i on i.benefit = b.benefit
     where section = $sec
--       and (($lang = 0) or (lang=$lang))
     order by b.section, secsort;
end
###
delimiter ;
-- call klpm_bensecitems_query(0);

delimiter ###
create or replace procedure klpm_benefits_family_query()
begin
    select x.family, x.benefit, x.lang, x.label
      from benl_family x
     order by x.family, x.benefit
    ;
end
###
delimiter ;


delimiter ###
create or replace procedure klpm_benefits_addon_query()
begin
    select x.addon, x.benefit, x.lang, x.label 
      from benl_addon x
     order by x.addon, x.benefit
    ;
end
###
delimiter ;
