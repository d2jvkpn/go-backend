CREATE TYPE account_status AS ENUM('created', 'activated', 'blocked', 'deleted');

CREATE TYPE account_level AS ENUM('admin', 'editor', 'reviewer', 'user', 'guest');
-- CREATE TYPE account_roles AS ENUM();

CREATE TABLE user_accounts (
  id          uuid DEFAULT gen_random_uuid(),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  status      account_status NOT NULL DEFAULT 'created',

  firstname  varchar(24) NOT NULL, -- UNIQUE,
  lastname   varchar(24) NOT NULL,
  phone      varchar(20) DEFAULT NULL UNIQUE,
  email      varchar(128) DEFAULT NULL UNIQUE,
  level      account_level NOT NULL,
  labels     varchar[] NOT NULL DEFAULT array[]::varchar[],

  password  varchar NOT NULL, -- bcrypt(password)

  PRIMARY KEY (id)
);

COMMENT ON TABLE user_accounts IS 'user accounts';
COMMENT ON COLUMN user_accounts.id IS 'account id';

CREATE TRIGGER updated_at BEFORE UPDATE ON user_accounts
  FOR EACH ROW EXECUTE PROCEDURE update_now();

CREATE INDEX user_accounts_created_at ON user_accounts (created_at DESC, status);
CREATE INDEX user_accounts_level ON user_accounts (level, created_at DESC);
CREATE INDEX user_accounts_labels ON user_accounts (labels);
