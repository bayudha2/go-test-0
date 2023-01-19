CREATE TABLE IF NOT EXISTS "public"."products" (
    "id" varchar(36) UNIQUE PRIMARY KEY NOT NULL,
    "name" varchar(36) NOT NULL,
    "price" numeric(10,2) NOT NULL DEFAULT 0.00,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);