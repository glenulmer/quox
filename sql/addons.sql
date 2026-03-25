delimiter ###
create or replace procedure quo_plan_categ_addons($plan int)
begin
    with
        xcross as (
            select p.plan, p.hospital, p.dental, px.segmask, c.categ
                 , c.required, a.product addon, a.level, x.famstate
              from plans p
              join products px on p.plan = px.product and !p.softdel and !px.softdel
              join families f on p.family = f.family and !f.softdel
              join family_products x on x.family = f.family and x.famstate > 0 and !x.softdel
              join products a on a.product = x.product and a.categ > 0 and a.display and !a.softdel
              join categs c on c.categ = a.categ and c.display
             where (a.segmask & px.segmask) <> 0
               and (($plan = 0) or (p.plan = $plan))
        ),
        none_options as (
            select distinct x.plan, x.segmask, x.categ
                 , 0 addon, 0 level, 0 famstate
              from xcross x
             where (x.categ = 3 and x.hospital = 30)
                or (x.categ = 4 and x.dental = 40)
        ),
        options as (
            select plan, segmask, categ, addon, level, famstate
              from xcross
            union all
            select plan, segmask, categ, addon, level, famstate
              from none_options
        )
        select o.plan, o.addon, o.categ
             , ifnull(a.level, 0) level, ifnull(l.label, ifnull(z.label, '---')) label
          from options o
          join plans p on p.plan = o.plan
          left join products a on o.addon = a.product
          left join levels l on a.level = l.level
          left join levels z on o.categ in (3, 4) and z.level = (10 * o.categ)
         where o.addon = 0 or l.level is not null
      order by o.plan, o.categ
          , if(o.categ in (3, 4), o.addon <> 0, o.addon = 0)
          , if(o.famstate between 1 and 3, 3 - o.famstate, if(o.famstate = 4, 3, 9))
          , o.level, o.addon;
end
###
delimiter ;
