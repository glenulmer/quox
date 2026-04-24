delimiter ###
create or replace procedure quo_today_get()
begin
    select convert(curdate(),int) today;
end
###
delimiter ;

create or replace view prices as
    select
        b.year as year,
        b.age as age,
        b.product as product,
        b.base as base,
        cast(ifnull (if (b.age between 21 and 59 and c.catsur <> 0,
                    round(b.base / 10, 0), 0), 0) as decimal(10,0)
        ) as surcharge
    from prices_base b
    join products p
    on b.product = p.product
    join categs c
    on c.categ = p.categ;

delimiter ###
create or replace procedure quo_segments_chooser()
begin
    select s.segment, s.name
      from segments s
     order by s.segment;
end
###
delimiter ;

delimiter ###
create or replace procedure quo_level_chooser_max($categ int, $ismax bool)
begin
    if $ismax then
        select level, label from levels l where l.categ = $categ order by level desc;
    else
        select level, label from levels l where l.categ = $categ order by level;
    end if;
end
###
delimiter ;


delimiter ###
create or replace procedure quo_segments_query()
begin
    select s.segment, s.name, s.code
      from segments s
  order by s.segment;
end
###
delimiter ;

delimiter ###
create or replace procedure quo_categs_query()
begin
    select c.categ, c.name, c.catsur, c.required
      from categs c
     where c.display = 1
  order by c.categ;
end
###
delimiter ;

alter table levels add column if not exists label varchar(16) after categ;
update levels set label = left(name,16) where label is null or label = '';
update levels set label = 'No Hospital' where level = 30;
update levels set label = 'No Dental' where level = 40;

delimiter ###
create or replace procedure quo_levels_query()
begin
  select l.level, l.label, c.categ, l.segments, (l.name != c.name) canStack
    from levels l
    join categs c on l.categ = c.categ
   where c.display = 1
  order by l.level;
end
###
delimiter ;

delimiter ###
create or replace function euro_whole($amount int)
returns varchar(12)
deterministic
begin
    return concat(replace(format(round($amount div 100, 0), 0), ',', '.'), ' €');
end
###
delimiter ;

delimiter ###
create or replace procedure quo_deductibles_chooser($adult bool, $ismax bool)
begin
    if $adult then
        if $ismax then
            select distinct ad_value div 100, euro_whole(ad_value) from plan_deductibles order by ad_value desc;
        else
            select distinct ad_value div 100, euro_whole(ad_value) from plan_deductibles order by ad_value;
        end if;
    else
        if $ismax then
            select distinct ch_value div 100, euro_whole(ch_value) from plan_deductibles order by ch_value desc;
        else
            select distinct ch_value div 100, euro_whole(ch_value) from plan_deductibles order by ch_value;
        end if;
    end if;
end
###
delimiter ;

update priorcov set descrip='No prior cover' where priorcov=0;

delimiter ###
create or replace procedure quo_priorcov_chooser()
begin
    select priorcov, descrip 
      from priorcov
  order by priorcov;
end
###
delimiter ;

delimiter ###
create or replace procedure quo_noexam_chooser()
begin
    select 0 noexam, 'Exam OK' descrip union all
    select 1, 'No exam';
end
###
delimiter ;

delimiter ###
create or replace procedure quo_specialist_chooser()
begin
    select 2 value, 'Not important' descrip union all
    select 1, 'Only referral' union all
    select 0, 'No referral';
end
###
delimiter ;

delimiter ###
create or replace procedure `quo_year_get`(in $year int)
begin
    select year, maxshare, cover, (cover*2) maxcover, ltccap 
      from years
     where (($year = 0) or (year = $year))
  order by year;
end
###
delimiter ;


delimiter ###
create or replace procedure quo_plan_details_query()
begin
    select p.plan, p.family
         , p.hospital, p.dental
         , p.priorcov, p.noexam, p.referral
         , p.tempvisa, p.surch, p.shi
         , p.vis_pct, p.vis_dec2
         , p.comonths
         , pd.ad_value, pd.ad_pct, pd.ch_value, pd.ch_pct
         , pn.promise, pn.note 
         , pn.ad_months, pn.ad_flat, pn.ch_months, pn.ch_flat
         , a.name, c.name, c.exact_age, a.segmask
         , ifnull(n.note,'') note, ifnull(n.style,'') style
      from plans p
      join families f on p.family = f.family
      join products a on p.plan = a.product
      join providers c on a.provider = c.provider
      join plan_deductibles pd on pd.plan = p.plan
      join plan_noclaims pn on pn.plan = p.plan
      left join plan_topnote n on n.plan = p.plan
     where !p.softdel and !a.softdel and !c.softdel
  order by c.name, f.name, pd.ad_value desc;
end
###
delimiter ;
--call quo_plan_details_query;


delimiter ###
create or replace procedure quo_product_query()
begin
    select p.product, c.provider, p.name, p.categ, p.level, p.segmask
      from products p
      join providers c on p.provider = c.provider
     where !p.softdel and !c.softdel and p.display;
end
###
delimiter ;

delimiter ###
create or replace procedure quo_product_prices()
begin
    select p.year, p.age, p.product, p.base, p.surcharge
      from years y
      join prices p on p.year = y.year
     where p.base > 0
     order by p.year, p.age, p.product;
end
###
delimiter ;




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
