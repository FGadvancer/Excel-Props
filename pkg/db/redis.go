package db

import (
	"Excel-Props/pkg/config"
	"context"
	go_redis "github.com/go-redis/redis/v8"
	"time"
)

func initRedis() *Redis {
	var r go_redis.UniversalClient
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		r = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		_, err := r.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	} else {
		r = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		_, err := r.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	}
	return NewRedis(r)
}

type Redis struct {
	RDB go_redis.UniversalClient
}

func NewRedis(RDB go_redis.UniversalClient) *Redis {
	return &Redis{RDB: RDB}
}

const (
	SheetID = "TEMPLATE_CODE:"
)

func (d *Redis) LockSheetID(sheetID string) error {
	key := SheetID + sheetID
	return d.RDB.SetNX(context.Background(), key, 1, time.Minute).Err()
}
func (d *Redis) UnLockSheetID(sheetID string) error {
	key := SheetID + sheetID
	return d.RDB.Del(context.Background(), key).Err()
}
