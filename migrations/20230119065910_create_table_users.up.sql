CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" varchar(36) UNIQUE PRIMARY KEY NOT NULL,
    "fullname" varchar(30) NOT NULL,
    "username" varchar(50) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW()
);