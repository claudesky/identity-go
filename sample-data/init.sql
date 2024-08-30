create table users (
    id uuid primary key,
    password text,
    name text constraint name_length check (char_length(name) <= 255),
    email text constraint email_length check (char_length(email) <= 255),
    email_verified_on timestamp,
    phone_number text constraint phone_number_length check (char_length(phone_number) <= 255),
    phone_number_verified_on timestamp
);

create table audiences (
    id uuid primary key,
    name text
);

create table token_families (
    id uuid primary key,
    sub uuid not null,
    last_issued uuid not null,
    created_at timestamp not null,
    last_issued_at timestamp not null
);

insert into users (id, name, email, password)
values (
    '0615b123-1a98-405b-bc44-6d41ad6a193c',
    'admin',
    'admin@example.org',
    '$2a$12$xrjwIS2d.hptiD/CEKKqxO5kVYjuWcWwxTNeXDT2bQJRlJweKGLu.'
)
