-- +migrate Up
CREATE TYPE USER_ROLE AS ENUM ('student', 'moderator', 'teacher', 'administrator', 'owner');

CREATE TABLE classes (
  id           UUID NOT NULL,
  name         TEXT NOT NULL,
  current_unit UUID
);

ALTER TABLE ONLY classes
  ADD CONSTRAINT classes_pkey PRIMARY KEY (id);

CREATE TABLE members (
  id       UUID                                     NOT NULL,
  user_id  UUID                                     NOT NULL,
  class_id UUID                                     NOT NULL,
  role     USER_ROLE DEFAULT 'student' :: USER_ROLE NOT NULL,
  FOREIGN KEY ("class_id") REFERENCES classes ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
ALTER TABLE ONLY members
  ADD CONSTRAINT members_pkey PRIMARY KEY (id);
CREATE INDEX members_class_id_idx
  ON members USING BTREE (class_id);
CREATE UNIQUE INDEX members_user_id_class_id_idx
  ON members USING BTREE (user_id, class_id);
CREATE INDEX members_user_id_idx
  ON members USING BTREE (user_id);

-- +migrate Down
DROP TABLE classes;
DROP TABLE members;
DROP TYPE USER_ROLE;