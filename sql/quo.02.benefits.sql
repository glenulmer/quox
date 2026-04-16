
delimiter ###
create or replace procedure quo_bensections_query()
begin
    select b.section, b.name
      from benefit_sections b
     order by b.section;
end
###
delimiter ;


delimiter ###
create or replace procedure quo_bensecitems_query($sec int)
begin
    select b.benefit, b.label, b.slim
      from benefits b
     where (($sec = 0) or (section = $sec))
     order by b.section, secsort;
end
###
delimiter ;


delimiter ###
create or replace procedure quo_benefits_family_query()
begin
    select x.benefit, x.family, x.descrip 
      from benefit_family_map x
     order by x.family, x.benefit
    ;
end
###
delimiter ;


delimiter ###
create or replace procedure quo_benefits_addon_query()
begin
    select x.benefit, x.product addon, x.descrip 
      from benefit_product_map x
     order by x.product, x.benefit
    ;
end
###
delimiter ;
