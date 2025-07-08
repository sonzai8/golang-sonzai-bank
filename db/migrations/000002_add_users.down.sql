-- Xóa ràng buộc foreign key giữa accounts.owner và users.username
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS owner_currency_key;

-- Xóa index unique trên bảng accounts
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_owner_fkey;

-- Xóa bảng users
DROP TABLE IF EXISTS users;