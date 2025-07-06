
create table novel(
    id serial PRIMARY KEY,
    name varchar,
    author varchar,
    category varchar,
    status varchar,
    intro varchar,
    cover_url varchar,
    url varchar
);

create table chapter(
    id int,
    novel_id int,
    seq int,
    title varchar,
    url varchar,
    content text
);

drop table novel;
drop table chapter;

delete from novel;
delete from chapter;

select * from novel;
select count(*) from novel;
select * from chapter;