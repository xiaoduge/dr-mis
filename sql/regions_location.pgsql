
drop table if exists regions_location;

create table regions_location(
    -- COMMENT '主键 ID'
    region_id serial primary key,
    --  COMMENT '地址全名'
    region_fullname text not null, 
    --  COMMENT '行政号码'
    region_code integer not null unique,
    --  COMMENT '地址拼音'
    region_pinyin text not null,
    --  COMMENT '地址名'
    region_name text not null,
    --  COMMENT '纬度'
    region_lat real not null,
    --  COMMENT '经度'
    region_lng real not null,
    region_cidx text not null,
    --  COMMENT '行政层级,省级0,市级1,区县2'
    region_level integer not null,
    --  COMMENT '归属'
    region_belongs integer not null
);