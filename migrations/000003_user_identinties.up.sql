create table if not exists todoapp.user_identities(
    id serial primary key,
    user_id integer not null references todoapp.users(id),
    provider text not null,
    subject text not null,
    email text not null,
    created_at timestamptz not null default now(),

    unique(provider, subject)
);