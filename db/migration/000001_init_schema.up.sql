create table if not exists posts (
    id serial primary key,
    title text not null,
    content text not null,
    created_at timestamptz not null default (now()),
    updated_at timestamptz not null default (now()) check (updated_at >= created_at)
);