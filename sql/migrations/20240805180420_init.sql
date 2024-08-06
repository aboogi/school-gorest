-- +goose Up
-- +goose StatementBegin
create type material_type as enum (
  'article', 'video', 'presentation'
);

create table school_material (
  material_id uuid primary key,
  name varchar(256) not null,
  content bytea null,
  material_type material_type not null,
  status bool not null,
  created_at timestamptz(0) not null,
  updated_at timestamptz(0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table school_material;
drop type material_type;
-- +goose StatementEnd
