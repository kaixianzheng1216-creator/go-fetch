-- +goose Up
create table if not exists users
(
    id            uuid primary key,
    username      varchar(255) not null unique,
    password_hash varchar(60)  not null,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz,
    deleted_at    timestamptz
);

comment on column users.id is 'ID';
comment on column users.username is '用户名';
comment on column users.password_hash is '密码哈希';
comment on column users.created_at is '创建时间';
comment on column users.updated_at is '更新时间';
comment on column users.deleted_at is '删除时间';

create table if not exists websites
(
    id         uuid primary key,
    user_id    uuid not null references users (id) on delete cascade,
    name       varchar(100) not null,
    domain     varchar(500),
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz
);

comment on column websites.id is 'ID';
comment on column websites.user_id is '用户';
comment on column websites.name is '名称';
comment on column websites.domain is '域名';
comment on column websites.created_at is '创建时间';
comment on column websites.updated_at is '更新时间';
comment on column websites.deleted_at is '删除时间';

create table if not exists sessions
(
    id          uuid primary key,
    website_id  uuid not null references websites (id) on delete cascade,
    browser     varchar(20),
    os          varchar(20),
    device      varchar(20),
    screen      varchar(11),
    language    varchar(35),
    country     char(2),
    region      varchar(20),
    city        varchar(50),
    distinct_id varchar(50),
    created_at  timestamptz not null default now()
);

comment on column sessions.id is 'ID';
comment on column sessions.website_id is '网站';
comment on column sessions.browser is '浏览器';
comment on column sessions.os is '操作系统';
comment on column sessions.device is '设备';
comment on column sessions.screen is '屏幕';
comment on column sessions.language is '语言';
comment on column sessions.country is '国家';
comment on column sessions.region is '地区';
comment on column sessions.city is '城市';
comment on column sessions.distinct_id is '访客标识';
comment on column sessions.created_at is '创建时间';

create table if not exists events
(
    id              uuid primary key,
    website_id      uuid         not null references websites (id) on delete cascade,
    session_id      uuid         not null references sessions (id) on delete cascade,
    visit_id        uuid         not null,
    event_type      integer      not null,
    event_name      varchar(50),
    url_path        varchar(500) not null,
    url_query       varchar(500),
    referrer_path   varchar(500),
    referrer_query  varchar(500),
    referrer_domain varchar(500),
    page_title      varchar(500),
    hostname        varchar(100),
    utm_source      varchar(255),
    utm_medium      varchar(255),
    utm_campaign    varchar(255),
    utm_content     varchar(255),
    utm_term        varchar(255),
    created_at      timestamptz not null default now()
);

comment on column events.id is 'ID';
comment on column events.website_id is '网站';
comment on column events.session_id is '会话';
comment on column events.visit_id is '访问';
comment on column events.event_type is '类型';
comment on column events.event_name is '事件名';
comment on column events.url_path is '路径';
comment on column events.url_query is '查询参数';
comment on column events.referrer_path is '来源路径';
comment on column events.referrer_query is '来源参数';
comment on column events.referrer_domain is '来源域名';
comment on column events.page_title is '页面标题';
comment on column events.hostname is '主机名';
comment on column events.utm_source is 'UTM 来源';
comment on column events.utm_medium is 'UTM 媒介';
comment on column events.utm_campaign is 'UTM 活动';
comment on column events.utm_content is 'UTM 内容';
comment on column events.utm_term is 'UTM 关键词';
comment on column events.created_at is '创建时间';

create table if not exists event_data
(
    id           uuid primary key,
    website_id   uuid         not null references websites (id) on delete cascade,
    event_id     uuid         not null references events (id) on delete cascade,
    data_key     varchar(500) not null,
    string_value varchar(500),
    number_value decimal(19, 4),
    date_value   timestamptz,
    data_type    integer      not null,
    created_at   timestamptz not null default now()
);

comment on column event_data.id is 'ID';
comment on column event_data.website_id is '网站';
comment on column event_data.event_id is '事件';
comment on column event_data.data_key is '键';
comment on column event_data.string_value is '字符串值';
comment on column event_data.number_value is '数值';
comment on column event_data.date_value is '日期值';
comment on column event_data.data_type is '数据类型';
comment on column event_data.created_at is '创建时间';

create table if not exists app_sessions
(
    token  text primary key,
    data   bytea       not null,
    expiry timestamptz not null
);

comment on column app_sessions.token is '令牌';
comment on column app_sessions.data is '数据';
comment on column app_sessions.expiry is '过期时间';

create index if not exists websites_user_idx on websites (user_id) where deleted_at is null;
create index if not exists events_website_type_created_idx on events (website_id, event_type, created_at);
create index if not exists events_website_session_created_idx on events (website_id, session_id, created_at);
create index if not exists events_website_visit_created_idx on events (website_id, visit_id, created_at);
create index if not exists events_website_path_created_idx on events (website_id, url_path, created_at);
create index if not exists events_website_event_created_idx on events (website_id, event_name, created_at);
create index if not exists event_data_event_idx on event_data (event_id);
create index if not exists app_sessions_expiry_idx on app_sessions (expiry);
