package helper

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

var ctx = context.Background()

func GetCount(db *gorm.DB, rdb *redis.Client, model any) (int64, error) {

	// 🔥 fallback (redis yo‘q bo‘lsa)
	if rdb == nil {
		var total int64
		if err := db.Model(model).Count(&total).Error; err != nil {
			return 0, err
		}
		return total, nil
	}

	// 🔥 table name olish
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return 0, err
	}

	table := stmt.Schema.Table
	key := "count:" + table

	// 🔹 Redisdan olish
	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		count, _ := strconv.ParseInt(val, 10, 64)
		return count, nil
	}

	// 🔹 DB dan olish
	var total int64
	if err := db.Model(model).Count(&total).Error; err != nil {
		return 0, err
	}

	// 🔹 Redisga yozish
	rdb.Set(ctx, key, total, time.Minute*5)

	return total, nil
}

func getModelType[T any](model *[]T) any {
	var entity T
	return &entity
}
