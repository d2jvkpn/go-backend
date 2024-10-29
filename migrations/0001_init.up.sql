CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE ok AS ENUM ('yes', 'no');

CREATE OR REPLACE FUNCTION update_now() RETURNS trigger AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$LANGUAGE plpgsql;
