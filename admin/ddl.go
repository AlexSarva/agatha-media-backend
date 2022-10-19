package admin

// ddl tables and queries for the first initializing of database
const ddl = `
CREATE TABLE if not exists public.users (
    id uuid primary key,
    username text,
    email text unique,
    passwd text,
    token text,
    token_expires timestamp,
    created timestamptz default now()
);
`
