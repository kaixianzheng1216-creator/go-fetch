-- name: CountUsers :one
select count(*)::bigint from users;

-- name: CreateUser :exec
insert into users (id, username, password_hash)
values (sqlc.arg(id)::uuid, sqlc.arg(username), sqlc.arg(password_hash));

-- name: GetUserByUsername :one
select id, username, password_hash, created_at
from users
where username = sqlc.arg(username);

-- name: GetUserByID :one
select id, username, password_hash, created_at
from users
where id = sqlc.arg(id)::uuid;

-- name: ListWebsites :many
select id, name, coalesce(domain, '')::text as domain, created_at
from websites
where user_id = sqlc.arg(user_id)::uuid and deleted_at is null
order by name;

-- name: CreateWebsite :one
insert into websites (id, user_id, name, domain)
values (sqlc.arg(id)::uuid, sqlc.arg(user_id)::uuid, sqlc.arg(name), nullif(sqlc.arg(domain)::text, ''))
returning id, name, coalesce(domain, '')::text as domain, created_at;

-- name: GetWebsite :one
select id, name, coalesce(domain, '')::text as domain, created_at
from websites
where id = sqlc.arg(id)::uuid and user_id = sqlc.arg(user_id)::uuid and deleted_at is null;

-- name: GetWebsiteForCollection :one
select id, name, coalesce(domain, '')::text as domain, created_at
from websites
where id = sqlc.arg(id)::uuid and deleted_at is null;

-- name: UpdateWebsite :execrows
update websites
set name = sqlc.arg(name), domain = nullif(sqlc.arg(domain)::text, ''), updated_at = now()
where id = sqlc.arg(id)::uuid and user_id = sqlc.arg(user_id)::uuid and deleted_at is null;

-- name: DeleteWebsite :execrows
update websites
set deleted_at = now()
where id = sqlc.arg(id)::uuid and user_id = sqlc.arg(user_id)::uuid and deleted_at is null;

-- name: InsertSession :exec
insert into sessions (
	id, website_id, browser, os, device, screen, language, country, region, city, distinct_id, created_at
)
values (
	sqlc.arg(id)::uuid, sqlc.arg(website_id)::uuid, nullif(sqlc.arg(browser)::text, ''),
	nullif(sqlc.arg(os)::text, ''), nullif(sqlc.arg(device)::text, ''), nullif(sqlc.arg(screen)::text, ''),
	nullif(sqlc.arg(language)::text, ''), nullif(sqlc.arg(country)::text, ''),
	nullif(sqlc.arg(region)::text, ''), nullif(sqlc.arg(city)::text, ''),
	nullif(sqlc.arg(distinct_id)::text, ''), sqlc.arg(created_at)
)
on conflict (id) do nothing;

-- name: InsertEvent :exec
insert into events (
	id, website_id, session_id, visit_id, event_type, event_name, url_path, url_query,
	referrer_path, referrer_query, referrer_domain, page_title, hostname, utm_source, utm_medium,
	utm_campaign, utm_content, utm_term, created_at
)
values (
	sqlc.arg(id)::uuid, sqlc.arg(website_id)::uuid, sqlc.arg(session_id)::uuid,
	sqlc.arg(visit_id)::uuid, sqlc.arg(event_type), nullif(sqlc.arg(event_name)::text, ''),
	sqlc.arg(url_path), nullif(sqlc.arg(url_query)::text, ''), nullif(sqlc.arg(referrer_path)::text, ''),
	nullif(sqlc.arg(referrer_query)::text, ''), nullif(sqlc.arg(referrer_domain)::text, ''),
	nullif(sqlc.arg(page_title)::text, ''), nullif(sqlc.arg(hostname)::text, ''), nullif(sqlc.arg(utm_source)::text, ''),
	nullif(sqlc.arg(utm_medium)::text, ''), nullif(sqlc.arg(utm_campaign)::text, ''),
	nullif(sqlc.arg(utm_content)::text, ''), nullif(sqlc.arg(utm_term)::text, ''),
	sqlc.arg(created_at)
);

-- name: InsertEventData :exec
insert into event_data (
	id, website_id, event_id, data_key, string_value, number_value, date_value, data_type, created_at
)
values (
	sqlc.arg(id)::uuid, sqlc.arg(website_id)::uuid, sqlc.arg(event_id)::uuid,
	sqlc.arg(data_key), nullif(sqlc.arg(string_value)::text, ''), sqlc.arg(number_value),
	sqlc.arg(date_value), sqlc.arg(data_type), sqlc.arg(created_at)
);

-- name: WebsiteStats :one
with visits as (
	select
		visit_id,
		min(session_id::text) as session_id,
		count(*) as pageviews,
		min(created_at) as min_time,
		max(created_at) as max_time
	from events
	where website_id = sqlc.arg(website_id)::uuid
	  and created_at between sqlc.arg(start_at) and sqlc.arg(end_at)
	  and event_type = sqlc.arg(pageview_event_type)
	group by visit_id
)
select
	coalesce(sum(pageviews), 0)::bigint as pageviews,
	count(distinct session_id)::bigint as visitors,
	count(*)::bigint as visits,
	count(*) filter (where pageviews = 1)::bigint as bounces,
	coalesce(sum(extract(epoch from (max_time - min_time))), 0)::bigint as total_time
from visits;

-- name: Pageviews :many
select
	date_trunc(sqlc.arg(bucket)::text, created_at)::timestamptz as time,
	count(*)::bigint as views,
	count(distinct session_id)::bigint as visitors
from events
where website_id = sqlc.arg(website_id)::uuid
  and created_at between sqlc.arg(start_at) and sqlc.arg(end_at)
  and event_type = sqlc.arg(pageview_event_type)
group by time
order by time;

-- name: EventMetrics :many
select
	coalesce(nullif(case sqlc.arg(metric)::text
		when 'path' then url_path
		when 'referrer' then referrer_domain
		when 'event' then event_name
		else ''
	end, ''), '(none)')::text as name,
	count(*)::bigint as views,
	count(distinct session_id)::bigint as visitors
from events
where website_id = sqlc.arg(website_id)::uuid
  and created_at between sqlc.arg(start_at) and sqlc.arg(end_at)
  and event_type = sqlc.arg(event_type)
group by name
order by views desc
limit sqlc.arg(limit_count);

-- name: SessionMetrics :many
select
	coalesce(nullif(case sqlc.arg(metric)::text
		when 'browser' then sessions.browser
		when 'os' then sessions.os
		when 'device' then sessions.device
		when 'country' then sessions.country
		else ''
	end, ''), '(none)')::text as name,
	count(*)::bigint as views,
	count(distinct events.session_id)::bigint as visitors
from events
join sessions on sessions.id = events.session_id and sessions.website_id = events.website_id
where events.website_id = sqlc.arg(website_id)::uuid
  and events.created_at between sqlc.arg(start_at) and sqlc.arg(end_at)
  and events.event_type = sqlc.arg(event_type)
group by name
order by views desc
limit sqlc.arg(limit_count);
