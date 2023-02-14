CREATE TABLE IF NOT EXISTS "public"."comments" (
    "id" varchar(36) UNIQUE PRIMARY KEY NOT NULL,
    "post_id" varchar(36) NOT NULL,
    "user_id" varchar(36) NOT NULL,
    "content" text NOT NULL,
    "parent_id" varchar(36),
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
)