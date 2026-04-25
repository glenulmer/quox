delimiter ###
create or replace procedure klpm_family_tips_query()
begin
    select x.family, x.tip
      from family_tips x
     order by family, pos
    ;
end
###
delimiter ;

delimiter ###
create or replace procedure klpm_plan_nccategs_query($plan int)
begin
    select categ
      from plan_noclaim_categs x
     where plan = $plan
    ;
end
###
delimiter ;

