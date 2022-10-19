select *
from gosZakupki.Company
         limit 100;

create database reestr_company;

drop table reestr_company.org;
create table reestr_company.org (
                                    ogrn String default '',
                                    inn String default '',
                                    max_num Int64 default 0,
                                    created DateTime64(3) default now()
)  engine = MergeTree
         ORDER BY (ogrn, inn);

drop table reestr_company.org_full;
create table reestr_company.org_full (
                                         ogrn String default '',
                                         crc32 String default '',
                                         reg_date String default '',
                                         min_num Int64 default 0,
                                         max_num Int64 default 0,
                                         end_date String default '',
                                         opf_id String default '',
                                         okved_id String default '',
                                         inn String default '',
                                         kpp String default '',
                                         short_name String default '',
                                         full_name String default '',
                                         email String default '',
                                         pfr String default '',
                                         fss String default '',
                                         capital String default '',
                                         kladr String default '',
                                         index_num String default '',
                                         region String default '',
                                         area String default '',
                                         city String default '',
                                         settlement String default '',
                                         street String default '',
                                         house String default '',
                                         corpus String default '',
                                         apartment String default '',
                                         created DateTime64(3) default now()
)  engine = MergeTree
       ORDER BY (ogrn, inn);

select ogrn, inn, kpp,
       replaceAll(case when short_name = '' then full_name else short_name end,'&quot;','"') short_name,
       replaceAll(full_name,'&quot;','"') full_name, reg_date, end_date, okved_id, capital,
       cast(region as Int8) region_id,
       replaceAll(case when area = '' then '' else area end ||
       case when city = '' then '' when area = '' then city else ', '||city end ||
       case when settlement = '' then '' when city = '' then settlement else ', '||settlement end ||
       case when street = '' then '' when settlement ='' then street else ', '||street end ||
       case when house = '' then '' else ', '||house end ||
       case when corpus in ('','-') then '' else ', '||corpus end ||
       case when apartment in ('','-') then '' else ', пом. '||apartment end,'&quot;','"')  address
       from reestr_company.org_full
        where inn = '9710091104';


select inn
from reestr_company.org_full
where inn = '9710091104';


insert into crawler.graphs (url_from, url_from_id, url_to, url_to_id, cnt_links)
with parents as (
    select distinct parent_url from crawler.posts where parent_url is not null
), parents_links as (
select url, base_url, srcs from crawler.posts
where 1=1
and url in (select parent_url from parents))
, raw_data as (
select posts.url url,
       posts.base_url main_url,
       posts.srcs,
       posts.parent_url parent,
       parents_links.base_url parent_base_url,
       parents_links.srcs parent_srcs
from crawler.posts
inner join parents_links on parents_links.url = posts.parent_url
where true
  and posts.parent_url is not null
  and posts.base_url != posts.parent_url
  and posts.url != posts.parent_url
  and posts.base_url != 'https://t.me/'
and posts.base_url != parents_links.base_url),
    url_ids_raw as (
        select distinct test from raw_data array join array(parent_base_url,main_url) as test
    ),
    url_ids as (
        select row_number() over (order by test) id, test from url_ids_raw
    ),
    main_data as (
        select raw_data.parent_base_url url_from, part1.id url_from_id,
               raw_data.main_url url_to, part2.id url_to_id
        from raw_data
                 left join url_ids part1 on part1.test = raw_data.parent_base_url
                 left join url_ids part2 on part2.test = raw_data.main_url
    )
select url_from, url_from_id, url_to, url_to_id, count() cnt_links
from main_data
group by url_from, url_from_id, url_to, url_to_id;

select * from crawler.graphs limit 10;





group by parent_base_url, main_url;








select url_from, url_to, cnt_links from crawler.graphs
    where 1=1
    and (url_from = 'https://t.me/moscowcurrent' or url_to = 'https://t.me/moscowcurrent');


with all_urls
select url_from from crawler.graphs
    union
select url_to from crawler.graphs;


drop table if exists crawler.graphs;
create table crawler.graphs
(
    url_from String,
    url_from_id Int64,
    url_to String,
    url_to_id Int64,
    cnt_links Int32,
    created   DateTime('Europe/Moscow') default now()
)
    engine = MergeTree()
        PARTITION BY (toYYYYMM(created))
        ORDER BY (url_from,url_to)
        SETTINGS index_granularity = 4000;

select * from crawler.posts where parent_url = 'https://t.me/a_beglov/932';

select count() from crawler.graphs
where graphs.cnt_links >= 10;


with sub_set as (
select test as url from crawler.graphs array join array(url_from,url_to) as test
where 1=1
  and (url_from = 'https://t.me/moscowach' or url_to = 'https://t.me/moscowach'))
select url_from,url_from_id, url_to, url_to_id, cnt_links from crawler.graphs
where 1=1
and (url_to in (select url from sub_set)
or url_from in (select url from sub_set))