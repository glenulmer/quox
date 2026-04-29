create or replace table languages
( lang int not null primary key
, label varchar(10) not null unique
, check (lang in (1,2))
, check (label in ('English','Deutsch'))
);

insert into languages values (1,'English');
insert into languages values (2,'Deutsch');

delimiter ###
create or replace procedure languages_query() 
begin
    select lang, label
      from languages x
     order by lang;
end
###
delimiter ;
call languages_query;

