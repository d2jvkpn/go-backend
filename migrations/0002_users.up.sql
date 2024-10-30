CREATE TYPE user_status AS ENUM('created', 'activated', 'blocked', 'deleted');

CREATE TYPE user_role AS ENUM('admin', 'manager', 'employee', 'customer', 'contractor');

CREATE TABLE user_accounts (
  id          uuid DEFAULT gen_random_uuid(),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  status      user_status NOT NULL,

  firstname  varchar(24) NOT NULL, -- UNIQUE,
  lastname   varchar(24) NOT NULL,
  phone      varchar(20) DEFAULT NULL UNIQUE,
  email      varchar(128) DEFAULT NULL UNIQUE,
  role       user_role NOT NULL,
  labels     varchar[] NOT NULL DEFAULT array[]::varchar[],

  password  varchar NOT NULL, -- bcrypt(password)

  PRIMARY KEY (id)
);

COMMENT ON TABLE user_accounts IS 'user accounts';
COMMENT ON COLUMN user_accounts.id IS 'account id';

CREATE TRIGGER updated_at BEFORE UPDATE ON user_accounts
  FOR EACH ROW EXECUTE PROCEDURE update_now();

CREATE INDEX user_accounts_created_at ON user_accounts (created_at DESC, status);
CREATE INDEX user_accounts_role ON user_accounts (role, created_at DESC);
CREATE INDEX user_accounts_labels ON user_accounts (labels);
