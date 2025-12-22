CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add the unique constraint to accounts
-- This ensures the combination of owner and currency is unique
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_currency_key" UNIQUE ("owner", "currency");

-- Link accounts to users
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");