-- +migrate Up
ALTER TYPE member_role RENAME TO user_role;
ALTER TYPE user_role ADD VALUE 'owner';
ALTER TABLE members ADD FOREIGN KEY ("class_id") REFERENCES classes("id") ON DELETE CASCADE ON UPDATE CASCADE;
