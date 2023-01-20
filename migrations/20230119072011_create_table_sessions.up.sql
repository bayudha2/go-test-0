CREATE TABLE IF NOT EXISTS "public"."sessions" (
    "id" varchar(36) UNIQUE PRIMARY KEY NOT NULL,
    "username" varchar(50) NOT NULL,
    "refresh_token" varchar NOT NULL,
    "expires_at" timestamptz NOT NULL DEFAULT NOW(),
    "created_at" timestamptz NOT NULL DEFAULT NOW()
);