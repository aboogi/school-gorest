-- name: GetSchoolMaterialByID :one
select * from school_material
where material_id = $1;

-- name: GetSchoolMaterialArray :many
select * from school_material
where 
  material_type = coalesce(sqlc.narg(material_type)::material_type, material_type)
  and case
    when sqlc.narg(from_created_at)::timestamptz is not null OR sqlc.narg(to_created_at)::timestamptz is not null
      then created_at between coalesce(@from_created_at, '2000-01-01 00:01:00+00') and coalesce(@to_created_at, now())
    else true end
order by created_at desc
limit @page_limit::bigint
offset @page_offset::bigint;

-- name: CreateSchoolMaterial :one
insert into school_material
(material_id, name, content, status, material_type, created_at)
values
($1, $2, $3, $4, $5, $6)
returning *;

-- name: UpdateSchoolMaterial :one
update school_material
set
  name = $2,
  content = $3,
  status = $4,
  updated_at = $5 
where material_id = $1
returning *;
