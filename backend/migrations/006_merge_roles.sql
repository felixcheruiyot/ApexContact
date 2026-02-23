-- Migration 006: Merge viewer / promoter / broadcaster → member
-- Drop the old CHECK constraint, remap existing data, re-add the new constraint.

-- 1. Drop the old role check constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- 2. Remap legacy roles to 'member'
UPDATE users SET role = 'member' WHERE role IN ('viewer', 'promoter', 'broadcaster');

-- 3. Add the new constraint
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('member', 'admin'));

-- 4. Update the column default
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'member';
