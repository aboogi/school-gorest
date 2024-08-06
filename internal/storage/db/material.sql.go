// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: material.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createSchoolMaterial = `-- name: CreateSchoolMaterial :one
insert into school_material
(material_id, name, content, status, material_type, created_at)
values
($1, $2, $3, $4, $5, $6)
returning material_id, name, content, material_type, status, created_at, updated_at
`

type CreateSchoolMaterialParams struct {
	MaterialID   uuid.UUID
	Name         string
	Content      []byte
	Status       bool
	MaterialType MaterialType
	CreatedAt    time.Time
}

func (q *Queries) CreateSchoolMaterial(ctx context.Context, arg CreateSchoolMaterialParams) (SchoolMaterial, error) {
	row := q.queryRow(ctx, q.createSchoolMaterialStmt, createSchoolMaterial,
		arg.MaterialID,
		arg.Name,
		arg.Content,
		arg.Status,
		arg.MaterialType,
		arg.CreatedAt,
	)
	var i SchoolMaterial
	err := row.Scan(
		&i.MaterialID,
		&i.Name,
		&i.Content,
		&i.MaterialType,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSchoolMaterialArray = `-- name: GetSchoolMaterialArray :many
select material_id, name, content, material_type, status, created_at, updated_at from school_material
where 
  material_type = coalesce($1::material_type, material_type)
  and case
    when $2::timestamptz is not null OR $3::timestamptz is not null
      then created_at between coalesce($2, '2000-01-01 00:01:00+00') and coalesce($3, now())
    else true end
order by created_at desc
limit $5::bigint
offset $4::bigint
`

type GetSchoolMaterialArrayParams struct {
	MaterialType  NullMaterialType
	FromCreatedAt sql.NullTime
	ToCreatedAt   sql.NullTime
	PageOffset    int64
	PageLimit     int64
}

func (q *Queries) GetSchoolMaterialArray(ctx context.Context, arg GetSchoolMaterialArrayParams) ([]SchoolMaterial, error) {
	rows, err := q.query(ctx, q.getSchoolMaterialArrayStmt, getSchoolMaterialArray,
		arg.MaterialType,
		arg.FromCreatedAt,
		arg.ToCreatedAt,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SchoolMaterial
	for rows.Next() {
		var i SchoolMaterial
		if err := rows.Scan(
			&i.MaterialID,
			&i.Name,
			&i.Content,
			&i.MaterialType,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSchoolMaterialByID = `-- name: GetSchoolMaterialByID :one
select material_id, name, content, material_type, status, created_at, updated_at from school_material
where material_id = $1
`

func (q *Queries) GetSchoolMaterialByID(ctx context.Context, materialID uuid.UUID) (SchoolMaterial, error) {
	row := q.queryRow(ctx, q.getSchoolMaterialByIDStmt, getSchoolMaterialByID, materialID)
	var i SchoolMaterial
	err := row.Scan(
		&i.MaterialID,
		&i.Name,
		&i.Content,
		&i.MaterialType,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateSchoolMaterial = `-- name: UpdateSchoolMaterial :one
update school_material
set
  name = $2,
  content = $3,
  status = $4,
  updated_at = $5 
where material_id = $1
returning material_id, name, content, material_type, status, created_at, updated_at
`

type UpdateSchoolMaterialParams struct {
	MaterialID uuid.UUID
	Name       string
	Content    []byte
	Status     bool
	UpdatedAt  sql.NullTime
}

func (q *Queries) UpdateSchoolMaterial(ctx context.Context, arg UpdateSchoolMaterialParams) (SchoolMaterial, error) {
	row := q.queryRow(ctx, q.updateSchoolMaterialStmt, updateSchoolMaterial,
		arg.MaterialID,
		arg.Name,
		arg.Content,
		arg.Status,
		arg.UpdatedAt,
	)
	var i SchoolMaterial
	err := row.Scan(
		&i.MaterialID,
		&i.Name,
		&i.Content,
		&i.MaterialType,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}