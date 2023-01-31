CREATE TABLE IF NOT EXISTS "public"."posts" (
    "id" varchar(36) UNIQUE PRIMARY KEY NOT NULL,
    "user_id" varchar(36) NOT NULL,
    "description" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);