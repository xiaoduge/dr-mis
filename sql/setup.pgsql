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
drop table if exists prizeinfo;
drop table if exists expiredprizeinfo;
drop table if exists userinfo;
create table userinfo(
    id serial primary key,
    user_id text not null unique,
    user_name text not null,
    user_code text not null unique,
    create_time timestamp not null
);
create table prizeinfo(
    prize_id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    user_code text references userinfo(user_code),
    prize_name text not null,
    prize_status integer not null, 
    created_time timestamp not null
);
create table expiredprizeinfo(
    prize_id serial primary key,
    user_id text references userinfo(user_id),
    user_name text not null,
    user_code text references userinfo(user_code),
    prize_name text not null,
    prize_status integer not null, 
    expired_time timestamp not null
)