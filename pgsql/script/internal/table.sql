CREATE SCHEMA test;

CREATE TABLE IF NOT EXISTS test.tb_aa (
  id SERIAL PRIMARY KEY,
  embed_aa Integer NOT NULL,
  embed_ab TEXT NOT NULL,
  note TEXT
);

CREATE TABLE IF NOT EXISTS test.tb_ab (
  pka Integer NOT NULL,
  pkb TEXT NOT NULL,
  embed_aa Integer NOT NULL,
  embed_ab TEXT NOT NULL,
  note TEXT,
  PRIMARY KEY (pka, pkb)
);

CREATE TABLE IF NOT EXISTS test.tb_array (
  i Integer[],
  t TEXT[]
);

