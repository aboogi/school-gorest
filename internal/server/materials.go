package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"school/internal/api"
	"school/internal/storage"
	"school/internal/storage/db"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Routes struct {
	storage *storage.Connection
	logger  *zap.Logger
}

// (GET /school/materials)
func (c Routes) GetSchoolMaterialArray(ctx echo.Context, params api.GetSchoolMaterialArrayParams) error {
	pagging := storage.Pagging(params.Page, params.PageSize)
	nmt, err := getNullMaterialType(params.MaterialType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse request: %w", err))
	}

	fromCreatedAt, er := getNullTime(params.FromCreatedAt)
	if er != nil {
		return er
	}

	toCreatedAt, er := getNullTime(params.ToCreatedAt)
	if er != nil {
		return er
	}

	rows, err := c.storage.Q.GetSchoolMaterialArray(
		ctx.Request().Context(), db.GetSchoolMaterialArrayParams{
			MaterialType:  nmt,
			FromCreatedAt: fromCreatedAt,
			ToCreatedAt:   toCreatedAt,
			PageOffset:    int64(pagging.Offset),
			PageLimit:     int64(pagging.Limit),
		},
	)
	if err != nil && err != sql.ErrNoRows {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse reques: %w", err))
	}

	res := api.SchoolMaterialArray{
		Materials: make([]api.SchoolMaterial, len(rows)),
	}

	for i := range rows {
		var status = api.SchoolMaterialStatusArchive
		if rows[i].Status {
			status = api.SchoolMaterialStatusActive
		}

		res.Materials[i] = api.SchoolMaterial{
			Content:      string(rows[i].Content),
			CreatedAt:    rows[i].CreatedAt,
			MaterialType: api.MaterialType(rows[i].MaterialType),
			Name:         rows[i].Name,
			Status:       status,
			UpdatedAt:    rows[i].UpdatedAt.Time,
			Id:           rows[i].MaterialID,
		}
	}
	return ctx.JSON(http.StatusOK, res)
}

// (POST /school/materials)
func (c Routes) CreateSchoolMaterial(ctx echo.Context) error {
	var req api.SchoolMaterial
	err := ctx.Bind(&req)
	if err != nil {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse request: %w", err))
	}

	status := true
	if req.Status == api.SchoolMaterialStatusArchive {
		status = false
	}

	nmt, err := getNullMaterialType(&req.MaterialType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse request: %w", err))
	}

	row, err := c.storage.Q.CreateSchoolMaterial(ctx.Request().Context(), db.CreateSchoolMaterialParams{
		MaterialID:   uuid.New(),
		Name:         req.Name,
		Content:      []byte(req.Content),
		Status:       status,
		MaterialType: nmt.MaterialType,
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse reques: %w", err))
	}

	var dbstatus = api.SchoolMaterialStatusArchive
	if row.Status {
		dbstatus = api.SchoolMaterialStatusActive
	}

	res := api.SchoolMaterial{
		Content:      string(row.Content),
		CreatedAt:    row.CreatedAt,
		MaterialType: api.MaterialType(row.MaterialType),
		Name:         row.Name,
		Status:       dbstatus,
		UpdatedAt:    row.UpdatedAt.Time,
		Id:           row.MaterialID,
	}

	return ctx.JSON(http.StatusCreated, res)
}

// (GET /school/materials/{uuid})
func (c Routes) GetSchoolMaterial(ctx echo.Context, materialID uuid.UUID) error {
	row, err := c.storage.Q.GetSchoolMaterialByID(ctx.Request().Context(), materialID)
	if err != nil {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse reques: %w", err))
	}

	var status = api.SchoolMaterialStatusArchive
	if row.Status {
		status = api.SchoolMaterialStatusActive
	}

	res := api.SchoolMaterial{
		Content:      string(row.Content),
		CreatedAt:    row.CreatedAt,
		MaterialType: api.MaterialType(row.MaterialType),
		Name:         row.Name,
		Status:       status,
		UpdatedAt:    row.UpdatedAt.Time,
		Id:           row.MaterialID,
	}
	return ctx.JSON(http.StatusOK, res)
}

// (PUT /school/materials/{uuid})
func (c Routes) UpdateSchoolMaterial(ctx echo.Context, materialID uuid.UUID) error {
	var req api.SchoolMaterial
	err := ctx.Bind(&req)
	if err != nil {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse reques: %w", err))
	}

	status := true
	if req.Status == api.SchoolMaterialStatusArchive {
		status = false
	}

	row, err := c.storage.Q.UpdateSchoolMaterial(ctx.Request().Context(), db.UpdateSchoolMaterialParams{
		MaterialID: materialID,
		Name:       req.Name,
		Content:    []byte(req.Content),
		Status:     status,
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		c.logger.Sugar().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error parse reques: %w", err))
	}

	var dbstatus = api.SchoolMaterialStatusArchive
	if row.Status {
		dbstatus = api.SchoolMaterialStatusActive
	}

	res := api.SchoolMaterial{
		Content:      string(row.Content),
		CreatedAt:    row.CreatedAt,
		MaterialType: api.MaterialType(row.MaterialType),
		Name:         row.Name,
		Status:       dbstatus,
		UpdatedAt:    row.UpdatedAt.Time,
		Id:           row.MaterialID,
	}

	return ctx.JSON(http.StatusOK, res)
}

func getNullMaterialType(materialType *api.MaterialType) (db.NullMaterialType, error) {
	if materialType == nil {
		return db.NullMaterialType{}, nil
	}

	var mt db.MaterialType
	switch *materialType {
	case api.Article:
		mt = db.MaterialTypeArticle
	case api.Video:
		mt = db.MaterialTypeVideo
	case api.Presentation:
		mt = db.MaterialTypePresentation
	default:
		return db.NullMaterialType{}, fmt.Errorf("unsupported value: %s", *materialType)
	}

	var nmt = db.NullMaterialType{
		MaterialType: mt,
		Valid:        true,
	}

	return nmt, nil
}

func getNullTime(t *time.Time) (sql.NullTime, error) {
	if t == nil {
		return sql.NullTime{}, nil
	}

	var nt sql.NullTime
	er := nt.Scan(*t)
	if er != nil {
		return sql.NullTime{}, er
	}

	return nt, nil
}
