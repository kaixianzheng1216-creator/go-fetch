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


create table if not exists event_data
(
    id           uuid primary key,
    website_id   uuid         not null references websites (id) on delete cascade,
    event_id     uuid         not null references events (id) on delete cascade,
    data_key     varchar(500) not null,
    string_value varchar(500),
    number_value double precision,
    date_value   timestamptz,
    data_type    integer      not null,
    created_at   timestamptz not null default now()
);


create table if not exists app_sessions
(
    token  text primary key,
    data   bytea       not null,
    expiry timestamptz not null
);

create index if not exists websites_user_active_idx
    on websites (user_id, deleted_at);

create index if not exists events_website_created_type_idx
    on events (website_id, created_at, event_type);

create index if not exists events_session_website_idx
    on events (session_id, website_id);

create index if not exists sessions_website_idx
    on sessions (website_id);

create index if not exists event_data_event_idx
    on event_data (event_id);

create index if not exists app_sessions_expiry_idx
    on app_sessions (expiry);
