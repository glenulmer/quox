alter table providers add column if not exists logo varchar(60) not null default '' after exact_age;
update providers set logo = concat(lower(name),'.jpg');
