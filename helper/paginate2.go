package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

///////////////////////
// REDIS
///////////////////////

var Rdb *redis.Client
var rCtx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

///////////////////////
// RESPONSE STRUCT
///////////////////////

type Meta struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}

type PaginatedResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

///////////////////////
// PAGINATE (MODEL)
///////////////////////

func Paginate[M any, R any](c echo.Context, db *gorm.DB, perPage int) (PaginatedResponse[R], error) {
	page, limit := parsePage(c, perPage)
	offset := (page - 1) * limit
	hasFilter := hasActiveFilter(c)

	var models []M
	var total int64

	if !hasFilter && Rdb != nil {
		dataKey := fmt.Sprintf("pg:%d:%d", page, limit)
		countKey := "pg:count"

		var res PaginatedResponse[R]

		if cached, err := Rdb.Get(rCtx, dataKey).Result(); err == nil {
			if json.Unmarshal([]byte(cached), &res) == nil {
				return res, nil
			}
		}

		if s, err := Rdb.Get(rCtx, countKey).Result(); err == nil {
			total, _ = strconv.ParseInt(s, 10, 64)
		} else {
			if err := countQuery[M](db).Count(&total).Error; err != nil {
				return res, err
			}
			Rdb.Set(rCtx, countKey, total, 5*time.Minute)
		}

		if err := db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
			return res, err
		}

		res = buildMappedResponse[M, R](models, total, page, limit)

		if raw, err := json.Marshal(res); err == nil {
			Rdb.Set(rCtx, dataKey, raw, 2*time.Minute)
		}

		return res, nil
	}

	if err := countQuery[M](db).Count(&total).Error; err != nil {
		return PaginatedResponse[R]{}, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return PaginatedResponse[R]{}, err
	}

	return buildMappedResponse[M, R](models, total, page, limit), nil
}

///////////////////////
// PAGINATE RAW (JOIN)
///////////////////////

func PaginateRaw[R any](c echo.Context, db *gorm.DB, perPage int) (PaginatedResponse[R], error) {
	page, limit := parsePage(c, perPage)
	offset := (page - 1) * limit
	hasFilter := hasActiveFilter(c)

	var total int64
	var data []R

	if !hasFilter && Rdb != nil {
		dataKey := fmt.Sprintf("pgraw:%d:%d", page, limit)
		countKey := "pgraw:count"

		var res PaginatedResponse[R]

		if cached, err := Rdb.Get(rCtx, dataKey).Result(); err == nil {
			if json.Unmarshal([]byte(cached), &res) == nil {
				return res, nil
			}
		}

		if s, err := Rdb.Get(rCtx, countKey).Result(); err == nil {
			total, _ = strconv.ParseInt(s, 10, 64)
		} else {
			if err := countRaw(db).Count(&total).Error; err != nil {
				return res, err
			}
			Rdb.Set(rCtx, countKey, total, 5*time.Minute)
		}

		if err := db.Offset(offset).Limit(limit).Scan(&data).Error; err != nil {
			return res, err
		}

		res = buildRawResponse(data, total, page, limit)

		if raw, err := json.Marshal(res); err == nil {
			Rdb.Set(rCtx, dataKey, raw, 2*time.Minute)
		}

		return res, nil
	}

	if err := countRaw(db).Count(&total).Error; err != nil {
		return PaginatedResponse[R]{}, err
	}

	if err := db.Offset(offset).Limit(limit).Scan(&data).Error; err != nil {
		return PaginatedResponse[R]{}, err
	}

	return buildRawResponse(data, total, page, limit), nil
}

///////////////////////
// COUNT
///////////////////////

func countQuery[M any](db *gorm.DB) *gorm.DB {
	var entity M
	stmt := db.Session(&gorm.Session{DryRun: true}).Find(&entity).Statement

	if stmt != nil && strings.Contains(strings.ToUpper(stmt.SQL.String()), " JOIN ") {
		return db.Session(&gorm.Session{NewDB: true}).
			Table("(?) as _c", gorm.Expr(stmt.SQL.String(), stmt.Vars...))
	}

	return db.Model(&entity)
}

func countRaw(db *gorm.DB) *gorm.DB {
	stmt := db.Session(&gorm.Session{DryRun: true}).Statement
	db.Find(nil)

	return db.Session(&gorm.Session{NewDB: true}).
		Table("(?) as _c", gorm.Expr(stmt.SQL.String(), stmt.Vars...))
}

///////////////////////
// BUILD RESPONSE
///////////////////////

func buildMappedResponse[M any, R any](models []M, total int64, page, limit int) PaginatedResponse[R] {
	data := make([]R, len(models))
	for i, m := range models {
		data[i] = AutoMap[M, R](m)
	}

	return PaginatedResponse[R]{
		Data: data,
		Meta: buildMeta(total, page, limit),
	}
}

func buildRawResponse[R any](data []R, total int64, page, limit int) PaginatedResponse[R] {
	return PaginatedResponse[R]{
		Data: data,
		Meta: buildMeta(total, page, limit),
	}
}

func buildMeta(total int64, page, limit int) Meta {
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	return Meta{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
		LastPage:    lastPage,
	}
}

///////////////////////
// AUTO MAP
///////////////////////

type convKind uint8

const (
	convAssign convKind = iota
	convConvert
)

type fieldStep struct {
	srcIdx int
	dstIdx int
	conv   convKind
}

type planKey struct{ src, dst reflect.Type }

var planCache sync.Map

func AutoMap[Src any, Dst any](src Src) Dst {
	var dst Dst
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(&dst).Elem()

	plan := cachedPlan(srcVal.Type(), dstVal.Type())

	for _, step := range plan {
		s := srcVal.Field(step.srcIdx)
		d := dstVal.Field(step.dstIdx)

		if step.conv == convConvert {
			d.Set(s.Convert(d.Type()))
		} else {
			d.Set(s)
		}
	}

	return dst
}

func cachedPlan(srcType, dstType reflect.Type) []fieldStep {
	key := planKey{srcType, dstType}

	if v, ok := planCache.Load(key); ok {
		return v.([]fieldStep)
	}

	dstIdx := make(map[string]int)

	for i := 0; i < dstType.NumField(); i++ {
		f := dstType.Field(i)
		name := strings.ToLower(f.Name)
		tag := f.Tag.Get("json")
		if tag != "" {
			name = strings.Split(tag, ",")[0]
		}
		dstIdx[name] = i
	}

	var plan []fieldStep

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		name := strings.ToLower(sf.Name)

		if idx, ok := dstIdx[name]; ok {
			df := dstType.Field(idx)

			step := fieldStep{srcIdx: i, dstIdx: idx}

			if sf.Type.ConvertibleTo(df.Type) {
				step.conv = convConvert
			}

			plan = append(plan, step)
		}
	}

	planCache.Store(key, plan)
	return plan
}

///////////////////////
// UTILS
///////////////////////

func parsePage(c echo.Context, perPage int) (page, limit int) {
	page, _ = strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ = strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = perPage
	}

	return
}

func hasActiveFilter(c echo.Context) bool {
	ignore := map[string]bool{
		"page": true, "limit": true, "column": true, "sort": true,
	}

	for key, vals := range c.QueryParams() {
		if ignore[key] {
			continue
		}
		for _, v := range vals {
			if v != "" {
				return true
			}
		}
	}
	return false
}
