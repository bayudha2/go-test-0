ALTER TABLE "public"."sessions"
    ADD CONSTRAINT "sessions_username_fkey" FOREIGN KEY ("username") REFERENCES "public"."users"("username");