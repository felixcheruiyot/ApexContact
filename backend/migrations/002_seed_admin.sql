-- Creates a default admin user if one doesn't already exist.
-- Password: Admin@1234  (bcrypt cost 10)
-- CHANGE THIS PASSWORD immediately after first login.
--
-- Run with:
--   docker compose exec postgres psql -U apex -d livestreamify -f /docker-entrypoint-initdb.d/002_seed_admin.sql
-- Or:
--   psql $DATABASE_URL -f migrations/002_seed_admin.sql

INSERT INTO users (id, email, password_hash, full_name, phone, role)
VALUES (
  uuid_generate_v4(),
  'admin@livestreamify.com',
  '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',  -- Admin@1234
  'Platform Admin',
  '',
  'admin'
)
ON CONFLICT (email) DO NOTHING;
