package infrastructure

import (
	"github.com/go-redis/redis"
	"taskmanage_api/src/constants"
)

var Redis = redis.NewClient(&redis.Options{
	Addr:     constants.Params.RedisAddr,
	Password: "", // no password set
	DB:       0,  // use default DB
})