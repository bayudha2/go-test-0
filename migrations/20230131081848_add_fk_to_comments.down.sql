ALTER TABLE "public"."comments"
    DROP CONSTRAINT "comments_post_id_fkey";

ALTER TABLE "public"."comments"
    DROP CONSTRAINT "comments_parent_id_fkey";

ALTER TABLE "public"."comments"
    DROP CONSTRAINT "comments_user_id_fkey";
