package dal

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"

	"permission-open/internal/pkg/conf"
)

var rdb *redis.ClusterClient

func MustInitRedis() {
	fmt.Println("redis initialized")
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(conf.Conf.Redis.ClusterAddress, ";"),
		Password: conf.Conf.Redis.Password,
	})
}

func GetRDB() *redis.ClusterClient {
	return rdb
}
