CREATE TABLE "verify_email"(
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secert_code" varchar NOT NULL,
    "is_used" bool NOT NULL DEFAULT false,
    "created_at" timestamptz not null default (now()),
    "expired_at" timestamptz not null default ((now() + interval '15 minutes'))
);

ALTER TABLE "verify_email" ADD FOREIGN KEY ("username") REFERENCES "users"("username");
ALTER TABLE "users" ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;