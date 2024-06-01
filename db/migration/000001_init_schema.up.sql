create table if not exists posts (
    id serial primary key,
    title text not null,
    content text not null,
    created_at timestamptz not null default (now()),
    updated_at timestamptz not null default (now()) check (updated_at >= created_at)
);

CREATE OR REPLACE FUNCTION update_modified_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_modified_time
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();
