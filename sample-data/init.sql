create table users (
    id uuid primary key,
    password text,
    name text constraint name_length check (char_length(name) <= 255),
    email text constraint email_length check (char_length(email) <= 255),
    email_verified_on timestamp,
    phone_number text constraint phone_number_length check (char_length(phone_number) <= 255),
    phone_number_verified_on timestamp
);

create table register_requests (
    id uuid primary key,
    password text not null,
    email text constraint email_length check (char_length(email) <= 255) not null,
    created_at timestamp not null
);

create table audiences (
    name text primary key
);

create table user_audience (
    user_id uuid not null,
    audience text not null,
    primary key (user_id, audience),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (audience) references audiences(name) on delete cascade
);

create table scopes (
    id uuid primary key,
    audience text not null,
    name text not null,
    foreign key (audience) references audiences(name) on delete cascade,
    unique (audience, name)
);

create table user_scope (
    user_id uuid not null,
    scope_id uuid not null,
    primary key (user_id, scope_id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (scope_id) references scopes(id) on delete cascade
);

create table apps (
    id uuid primary key,
    name text not null
);

create table token_families (
    id uuid primary key,
    sub uuid not null,
    last_issued uuid not null,
    created_at timestamp not null,
    last_issued_at timestamp not null,
    expires_at timestamp not null,
    revoked boolean not null default false
);

create table email_verification_requests (
    id uuid primary key,
    register_request_id uuid not null,
    email text not null,
    token text not null,
    revoked boolean not null default false,
    accepted boolean not null default false,
    expires_at timestamp not null,
    created_at timestamp not null
);

-- Admin User

insert into users (id, name, email, password)
values (
    '0615b123-1a98-405b-bc44-6d41ad6a193c',
    'admin',
    'admin@example.org',
    -- password is "password"
    '$2a$12$xrjwIS2d.hptiD/CEKKqxO5kVYjuWcWwxTNeXDT2bQJRlJweKGLu.'
);

insert into audiences (name)
values ('http://localhost:9102'); -- idg_app_url

insert into scopes (id, audience, name)
values (
    '0615b123-1a98-405b-bc44-6d41ad6a193d',
    'http://localhost:9102',
    'read'
);

insert into user_scope (user_id, scope_id)
values (
    '0615b123-1a98-405b-bc44-6d41ad6a193c',
    '0615b123-1a98-405b-bc44-6d41ad6a193d'
);
