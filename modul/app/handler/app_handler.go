package app_handler

import (
	"log"
	app_service "my-project/modul/app/service"

	"github.com/Abdugoffor/echo-crud-pg/request"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service app_service.AppService
}

func NewAppHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) {
	handler := AppHandler{
		db:      db,
		log:     log,
		service: app_service.NewAppService(db),
	}

	route := gorm.Group("/app")
	{
		route.GET("", handler.All)
		route.GET("/page", handler.Page)
		route.GET("/:id", handler.Show)
	}
}

func (handler *AppHandler) All(ctx echo.Context) error {
	req := request.Request(ctx)

	filter := func(tx *gorm.DB) *gorm.DB {

		tx.Select(`
		apps.id,
		to_char(apps.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
		to_char(apps.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at,
		to_char(apps.deleted_at, 'YYYY-MM-DD HH24:MI') AS deleted_at,
			jsonb_build_object(
			'id', app_categories.id, 
			'name', app_categories.name, 
			'slug', app_categories.slug
			) as app_category
		`).Joins("JOIN app_categories ON app_categories.id = apps.app_category_id")

		return tx
	}

	data, err := handler.service.All(req.Context(), req.NewPaginate(), filter)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(200, data)
}

func (handler *AppHandler) Page(ctx echo.Context) error {
	req := request.Request(ctx)

	filter := func(tx *gorm.DB) *gorm.DB {
		return tx.
			Select(`
				apps.id,
				apps.user_id,
				apps.is_active,

				to_char(apps.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
				to_char(apps.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at,
				to_char(apps.deleted_at, 'YYYY-MM-DD HH24:MI') AS deleted_at,

				-- CATEGORY
				jsonb_build_object(
					'id', ac.id,
					'name', ac.name,
					'slug', ac.slug
				) AS app_category,

				-- PAGES + FORMS + VALUES
				COALESCE(
					jsonb_agg(
						DISTINCT jsonb_build_object(
							'id', ap.id,
							'name', ap.name,
							'slug', ap.slug,
							'forms',
							COALESCE(
								(
									SELECT jsonb_agg(
										jsonb_build_object(
											'id', af2.id,
											'name', af2.name,
											'slug', af2.slug,
											'type', af2.type,
											'is_require', af2.is_require,
											'value', av2.value
										)
									)
									FROM app_form af2
									LEFT JOIN app_values av2
										ON av2.app_form_id = af2.id
										AND av2.app_id = apps.id
										AND av2.deleted_at IS NULL
									WHERE af2.app_page_id = ap.id
										AND af2.deleted_at IS NULL
								),
								'[]'::jsonb
							)
						)
					) FILTER (WHERE ap.id IS NOT NULL),
					'[]'::jsonb
				) AS app_page,

				-- nechta page bor (ixtiyoriy)
				COALESCE(
					COUNT(DISTINCT ap.id),
					0
				) AS page_count
			`).
			Joins(`JOIN app_categories ac ON ac.id = apps.app_category_id`).
			// categoryga tegishli barcha pagelar
			Joins(`LEFT JOIN app_pages ap ON ap.app_category_id = ac.id AND ap.deleted_at IS NULL`).
			// asosiy filterlar
			Where(`apps.deleted_at IS NULL`).
			// // bitta ariza bo‘yicha filter (agar ID querydan kelsa)
			// Where(`apps.id = ?`, ctx.Param("id")). // yoki req.Id bo‘lsa: Where("apps.id = ?", req.Id)
			Group(`
				apps.id,
				apps.user_id,
				apps.is_active,
				apps.created_at,
				apps.updated_at,
				apps.deleted_at,
				ac.id, ac.name, ac.slug
			`)
	}

	data, err := handler.service.Page(req.Context(), filter)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(200, data)
}

func (handler *AppHandler) Show(ctx echo.Context) error {
	req := request.Request(ctx) // bu yerda ID va boshqalar bo‘lsa ham bo‘ladi

	filter := func(tx *gorm.DB) *gorm.DB {
		return tx.
			Select(`
				apps.id,
				apps.user_id,
				apps.is_active,

				to_char(apps.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
				to_char(apps.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at,
				to_char(apps.deleted_at, 'YYYY-MM-DD HH24:MI') AS deleted_at,

				-- CATEGORY
				jsonb_build_object(
					'id', ac.id,
					'name', ac.name,
					'slug', ac.slug
				) AS app_category,

				-- PAGES + FORMS + VALUES
				COALESCE(
					jsonb_agg(
						DISTINCT jsonb_build_object(
							'id', ap.id,
							'name', ap.name,
							'slug', ap.slug,
							'forms',
							COALESCE(
								(
									SELECT jsonb_agg(
										jsonb_build_object(
											'id', af2.id,
											'name', af2.name,
											'slug', af2.slug,
											'type', af2.type,
											'is_require', af2.is_require,
											'value', av2.value
										)
									)
									FROM app_form af2
									LEFT JOIN app_values av2
										ON av2.app_form_id = af2.id
										AND av2.app_id = apps.id
										AND av2.deleted_at IS NULL
									WHERE af2.app_page_id = ap.id
										AND af2.deleted_at IS NULL
								),
								'[]'::jsonb
							)
						)
					) FILTER (WHERE ap.id IS NOT NULL),
					'[]'::jsonb
				) AS app_page,

				-- nechta page bor (ixtiyoriy)
				COALESCE(
					COUNT(DISTINCT ap.id),
					0
				) AS page_count
			`).
			Joins(`JOIN app_categories ac ON ac.id = apps.app_category_id`).
			// categoryga tegishli barcha pagelar
			Joins(`LEFT JOIN app_pages ap ON ap.app_category_id = ac.id AND ap.deleted_at IS NULL`).
			// asosiy filterlar
			Where(`apps.deleted_at IS NULL`).
			// bitta ariza bo‘yicha filter (agar ID querydan kelsa)
			Where(`apps.id = ?`, ctx.Param("id")). // yoki req.Id bo‘lsa: Where("apps.id = ?", req.Id)
			Group(`
				apps.id,
				apps.user_id,
				apps.is_active,
				apps.created_at,
				apps.updated_at,
				apps.deleted_at,
				ac.id, ac.name, ac.slug
			`)
	}

	data, err := handler.service.Show(req.Context(), filter)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(200, data)
}
