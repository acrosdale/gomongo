CREATE TABLE polls (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  question varchar(255) NOT NULL,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

--bun:split

CREATE TABLE comments (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  body text NOT NULL,
  poll_id BIGINT NOT NULL REFERENCES polls (id) ON DELETE CASCADE,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX comments_poll_id_idx ON comments(poll_id);
