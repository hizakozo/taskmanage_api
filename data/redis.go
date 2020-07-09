package data

import (
	"github.com/go-redis/redis"
	"fmt"
	"context"
	"encoding/json"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func RedisSet(json string, key string) {
	
	err := client.Set(ctx, key, json, 0).Err()
    if err != nil {
        fmt.Println("redis.Client.Set Error:", err)
    }
}

func RedisGet(key string) (User, error) {
	userInfoJson, _ := client.Get(ctx, key).Result()
	var user = new(User)
	err := json.Unmarshal([]byte(userInfoJson), user)
	var userInfo = User{user.ID, user.Name, user.Avatar, user.Isdelete}
	return userInfo, err
}

func RedisGetInviteInfo(token string) (string, error) {
	inviteInfoJson, err := client.Get(ctx, token).Result()
	return inviteInfoJson, err
}