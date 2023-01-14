start batch ddl;

CREATE TABLE if not exists users (
  user_id varchar not null primary key,
  name varchar not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

CREATE TABLE if not exists user_items (
  user_id varchar not null,
  item_id varchar not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  primary key(user_id, item_id)
) interleave in parent users on delete cascade;

create table if not exists items (
  item_id varchar not null primary key,
  item_name varchar not null,
  price numeric not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

run batch;
