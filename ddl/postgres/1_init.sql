-- +migrate Up
CREATE TYPE member_role AS ENUM('student','moderator','teacher','administrator');

CREATE TABLE classes (
  id uuid NOT NULL,
  name text NOT NULL,
  current_unit uuid
);

ALTER TABLE ONLY classes
  ADD CONSTRAINT classes_pkey PRIMARY KEY (id);

CREATE TABLE members (
  id uuid NOT NULL,
  user_id uuid NOT NULL,
  class_id uuid NOT NULL,
  role member_role DEFAULT 'student'::member_role NOT NULL
);
ALTER TABLE ONLY members
  ADD CONSTRAINT members_pkey PRIMARY KEY (id);
CREATE INDEX members_class_id_idx ON members USING btree (class_id);
CREATE UNIQUE INDEX members_user_id_class_id_idx ON members USING btree (user_id, class_id);
CREATE INDEX members_user_id_idx ON members USING btree (user_id);

-- +migrate Down
DROP TABLE classes;
DROP TABLE members;