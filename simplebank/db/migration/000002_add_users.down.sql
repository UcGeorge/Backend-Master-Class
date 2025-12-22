-- Remove the foreign key constraint
-- (Postgres defaults the name to table_column_fkey unless specified)
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Remove the unique constraint
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_currency_key";

-- Drop the users table
DROP TABLE IF EXISTS "users";