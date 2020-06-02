-- 门店信息
drop table if exists storeinfo;
create table storeinfo(
    store_id serial primary key,
    store_name text not null,
    store_phone text not null,

    store_province text not null,
    store_province_code integer not null,

    store_city text not null,
    store_city_code integer not null,

    store_county text,
    store_county_code integer,
    
    store_address text not null,
    store_tag text not null,
    store_img text not null unique,

    created_time timestamp not null
)

-- 领奖信息
drop table if exists playinfo;
drop table if exists gamerecord;
drop table if exists awardinfo;
drop table if exists awardedinfo;
drop table if exists userinfo;
drop table if exists prizesinfo;

-- 奖品信息记录表
create table prizesinfo(
    id serial       primary key,
    prize_name      text not null,
    prize_quantity  integer not null,
    prize_remaining integer not null,
    prize_used      integer not null,
    prize_category  text not null,
    mark            text not null,
    create_time     timestamp not null
);

-- 用户信息记录表
create table userinfo(
    id serial primary key,
    user_id text not null unique,
    user_name text not null,
    user_code text not null unique,
    mark text not null,
    create_time timestamp not null
);



-- 待领取的奖品信息表
create table awardinfo(
    id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    user_code text references userinfo(user_code),
    prize_name text not null,
    prize_id integer references prizesinfo(id),
    prize_status integer not null, 
    prize_category text not null,
    created_time timestamp not null
);

-- 已经领取的奖品信息表
create table awardedinfo(
    id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    user_code text references userinfo(user_code),
    prize_name text not null,
    prize_id integer references prizesinfo(id),
    prize_status integer not null, 
    prize_category text not null,
    expired_time timestamp not null
);

-- 游戏情况记录表
create table gamerecord(
    id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    score integer not null,
    spending integer not null,
    medal integer not null,
    mark text not null,
    update_time timestamp not null
);

-- 参与情况记录
create table playinfo(
    id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    play_times integer not null,
    draw_times integer not null,
    awards_times integer not null,
    canwin integer not null,
    mark text not null,
    update_time timestamp not null
);

-- 测试sql
select * from userinfo;
select * from gamerecord;
select * from playinfo;
select * from prizesinfo;
select * from awardinfo;
select * from awardedinfo;
