package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

var client = redis.NewClient(&redis.Options{
	//Addr:     "taskmanage-redis:6379",
	Addr:     "127.0.0.1:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func RedisSet(json string, key string) {
	
	err := client.Set(key, json, 0).Err()
    if err != nil {
        fmt.Println("redis.Client.Set Error:", err)
    }
}

func RedisGet(key string) (User, error) {
	userInfoJson, _ := client.Get(key).Result()
	var user = new(User)
	err := json.Unmarshal([]byte(userInfoJson), user)
	var userInfo = User{user.ID, user.Name, user.Avatar, user.Isdelete}
	return userInfo, err
}

func RedisGetInviteInfo(token string) (string, error) {
	inviteInfoJson, err := client.Get(token).Result()
	return inviteInfoJson, err
}

func RedisDelete(token string) {
	client.Del(token)
}