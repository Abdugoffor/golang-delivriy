package app_handler

import (
	"log"
	app_service "my-project/modul/app/service"
	"net/http"
	"strconv"

	"github.com/Abdugoffor/echo-crud-pg/request"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CateHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service app_service.AppCate
}

func NewAppCateHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) {
	handler := CateHandler{
		db:      db,
		log:     log,
		service: app_service.NewAppCate(db),
	}

	route := gorm.Group("/app/cate")
	{
		route.GET("", handler.All)
		route.GET("/:id", handler.Show)
	}
}

func (handler *CateHandler) All(ctx echo.Context) error {
	req := request.Request(ctx)

	filter := func(tx *gorm.DB) *gorm.DB {
		return tx.
			Select(`
			
			app_categories.id AS id,
			app_categories.name AS name,
			app_categories.slug AS slug,
			app_categories.is_active AS is_active,
			to_char(app_categories.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
			to_char(app_categories.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at,
			to_char(app_categories.deleted_at, 'YYYY-MM-DD HH24:MI') AS deleted_at,

			COALESCE(
				jsonb_agg(
					jsonb_build_object(
						'id', ap.id,
						'name', ap.name,
						'slug', ap.slug,
						'is_active', ap.is_active,
						'forms',
						COALESCE(
							(
								SELECT jsonb_agg(
									jsonb_build_object(
										'id', af.id,
										'name', af.name,
										'slug', af.slug,
										'type', af.type,
										'is_require', af.is_require,
										'is_active', af.is_active,
										'options',
										CASE
											WHEN af.type = 'select' THEN
												COALESCE(
													(
														SELECT jsonb_agg(
															jsonb_build_object(
																'id', ao.id,
																'name', ao.name,
																'slug', ao.slug,
																'is_active', ao.is_active
															)
														)
														FROM app_option ao
														WHERE ao.app_form_id = af.id
															AND ao.is_active = TRUE
															AND ao.deleted_at IS NULL
													),
													'[]'::jsonb
												)
											ELSE '[]'::jsonb
										END
									)
									ORDER BY af.id
								)
								FROM app_form af
								WHERE af.app_page_id = ap.id
									AND af.is_active = TRUE       
									AND af.deleted_at IS NULL
							),
							'[]'::jsonb
						)
					)
					ORDER BY ap.id      
				) FILTER (WHERE ap.id IS NOT NULL),
				'[]'::jsonb
			) AS pages
		`).
			Joins(`
			LEFT JOIN app_pages ap
				ON ap.app_category_id = app_categories.id
				AND ap.deleted_at IS NULL
		`).
			Where(`app_categories.deleted_at IS NULL`).
			Group(`
			app_categories.id,
			app_categories.name,
			app_categories.slug,
			app_categories.is_active,
			app_categories.created_at,
			app_categories.updated_at,
			app_categories.deleted_at
		`)
	}

	data, err := handler.service.All(req.Context(), req.NewPaginate(), filter)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(200, data)
}

func (handler *CateHandler) Show(ctx echo.Context) error {
	req := request.Request(ctx)
	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseInt(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.
			Select(`
			app_categories.id AS id,
			app_categories.name AS name,
			app_categories.slug AS slug,
			app_categories.is_active AS is_active,
			to_char(app_categories.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
			to_char(app_categories.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at,
			to_char(app_categories.deleted_at, 'YYYY-MM-DD HH24:MI') AS deleted_at,

			COALESCE(
				jsonb_agg(
					jsonb_build_object(
						'id', ap.id,
						'name', ap.name,
						'slug', ap.slug,
						'is_active', ap.is_active,
						'forms',
						COALESCE(
							(
								SELECT jsonb_agg(
									jsonb_build_object(
										'id', af.id,
										'name', af.name,
										'slug', af.slug,
										'type', af.type,
										'is_require', af.is_require,
										'is_active', af.is_active,
										'options',
										CASE
											WHEN af.type = 'select' THEN
												COALESCE(
													(
														SELECT jsonb_agg(
															jsonb_build_object(
																'id', ao.id,
																'name', ao.name,
																'slug', ao.slug,
																'is_active', ao.is_active
															)
														)
														FROM app_option ao
														WHERE ao.app_form_id = af.id
															AND ao.is_active = TRUE
															AND ao.deleted_at IS NULL
													),
													'[]'::jsonb
												)
											ELSE '[]'::jsonb
										END
									)
									ORDER BY af.id
								)
								FROM app_form af
								WHERE af.app_page_id = ap.id
									AND af.is_active = TRUE       
									AND af.deleted_at IS NULL
							),
							'[]'::jsonb
						)
					)
					ORDER BY ap.id      
				) FILTER (WHERE ap.id IS NOT NULL),
				'[]'::jsonb
			) AS pages
		`).
			Joins(`
			LEFT JOIN app_pages ap
				ON ap.app_category_id = app_categories.id
				AND ap.deleted_at IS NULL
		`).
			Where(`app_categories.deleted_at IS NULL`).
			Group(`
			app_categories.id,
			app_categories.name,
			app_categories.slug,
			app_categories.is_active,
			app_categories.created_at,
			app_categories.updated_at,
			app_categories.deleted_at
		`)

		if parsedID > 0 {
			tx = tx.Where("app_categories.id = ?", parsedID)
		}

		return tx
	}

	data, err := handler.service.Show(req.Context(), filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}
	return ctx.JSON(200, data)
}
