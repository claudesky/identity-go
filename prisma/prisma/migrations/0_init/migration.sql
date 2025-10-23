-- CreateTable
CREATE TABLE "apps" (
    "id" UUID NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "apps_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "audiences" (
    "name" TEXT NOT NULL,

    CONSTRAINT "audiences_pkey" PRIMARY KEY ("name")
);

-- CreateTable
CREATE TABLE "email_verification_requests" (
    "id" UUID NOT NULL,
    "register_request_id" UUID NOT NULL,
    "email" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "revoked" BOOLEAN NOT NULL DEFAULT false,
    "accepted" BOOLEAN NOT NULL DEFAULT false,
    "expires_at" TIMESTAMP(6) NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL,

    CONSTRAINT "email_verification_requests_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "register_requests" (
    "id" UUID NOT NULL,
    "password" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL,

    CONSTRAINT "register_requests_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "scopes" (
    "id" UUID NOT NULL,
    "audience" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "scopes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "token_families" (
    "id" UUID NOT NULL,
    "sub" UUID NOT NULL,
    "last_issued" UUID NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL,
    "last_issued_at" TIMESTAMP(6) NOT NULL,
    "expires_at" TIMESTAMP(6) NOT NULL,
    "revoked" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "token_families_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_audience" (
    "user_id" UUID NOT NULL,
    "audience" TEXT NOT NULL,

    CONSTRAINT "user_audience_pkey" PRIMARY KEY ("user_id","audience")
);

-- CreateTable
CREATE TABLE "user_scope" (
    "user_id" UUID NOT NULL,
    "scope_id" UUID NOT NULL,

    CONSTRAINT "user_scope_pkey" PRIMARY KEY ("user_id","scope_id")
);

-- CreateTable
CREATE TABLE "users" (
    "id" UUID NOT NULL,
    "password" TEXT,
    "name" TEXT,
    "email" TEXT,
    "email_verified_on" TIMESTAMP(6),
    "phone_number" TEXT,
    "phone_number_verified_on" TIMESTAMP(6),

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "scopes_audience_name_key" ON "scopes"("audience", "name");

-- AddForeignKey
ALTER TABLE "scopes" ADD CONSTRAINT "scopes_audience_fkey" FOREIGN KEY ("audience") REFERENCES "audiences"("name") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_audience" ADD CONSTRAINT "user_audience_audience_fkey" FOREIGN KEY ("audience") REFERENCES "audiences"("name") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_audience" ADD CONSTRAINT "user_audience_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_scope" ADD CONSTRAINT "user_scope_scope_id_fkey" FOREIGN KEY ("scope_id") REFERENCES "scopes"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_scope" ADD CONSTRAINT "user_scope_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

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
