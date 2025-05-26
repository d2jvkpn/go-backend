CREATE TYPE account_level AS ENUM('admin', 'editor', 'reviewer', 'user', 'guest');
-- CREATE TYPE account_roles AS ENUM();

CREATE TYPE account_status AS ENUM('created', 'activated', 'blocked', 'deleted');


CREATE TABLE user_accounts (
  id          uuid DEFAULT gen_random_uuid(),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  status      account_status NOT NULL DEFAULT 'created',

  firstname  varchar(32) NOT NULL, -- UNIQUE,
  lastname   varchar(32) NOT NULL,
  -- length range=[5, 64]
  email      varchar(64) DEFAULT NULL UNIQUE,
  phone      varchar(20) DEFAULT NULL UNIQUE,
  level      account_level NOT NULL,
  -- max=16x32
  labels     varchar[] NOT NULL DEFAULT array[]::varchar[],
  -- raw_password: 8-32, [a-z][A-Z][0-9]
  password  varchar NOT NULL, -- bcrypt(password)

  PRIMARY KEY (id)
);

COMMENT ON TABLE user_accounts IS 'user accounts';
COMMENT ON COLUMN user_accounts.id IS 'account id';

CREATE TRIGGER updated_at BEFORE UPDATE ON user_accounts
  FOR EACH ROW EXECUTE PROCEDURE update_now();

CREATE INDEX user_accounts_firstname_trgm ON user_accounts USING gin (firstname gin_trgm_ops);
CREATE INDEX user_accounts_lastname_trgm ON user_accounts USING gin (lastname gin_trgm_ops);

CREATE INDEX user_accounts_created_at ON user_accounts (created_at DESC, status);
CREATE INDEX user_accounts_level ON user_accounts (level, created_at DESC);
CREATE INDEX user_accounts_labels ON user_accounts (labels);
