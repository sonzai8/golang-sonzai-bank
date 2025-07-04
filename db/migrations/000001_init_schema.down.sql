-- Xoá ràng buộc trước
ALTER TABLE "entries" DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";

-- Xoá index
DROP INDEX IF EXISTS "accounts_owner_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_idx";
DROP INDEX IF EXISTS "transfers_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_to_account_id_idx";

-- Xoá bảng theo thứ tự phụ thuộc
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "accounts";