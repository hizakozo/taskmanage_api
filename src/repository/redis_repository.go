package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"taskmanage_api/src/domain"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.RedisRepository {
	return &redisRepository{
		client: client,
	}
}

func (rr *redisRepository) RedisSet(json string, key string) {

	err := rr.client.Set(key, json, 0).Err()
	if err != nil {
		fmt.Println("redis.Client.Set Error:", err)
	}
}

func (rr *redisRepository) RedisGet(key string) (domain.User, error) {
	userInfoJson, _ := rr.client.Get(key).Result()
	var user = new(domain.User)
	err := json.Unmarshal([]byte(userInfoJson), user)
	var userInfo = domain.User{user.ID, user.Name, user.Avatar, user.Isdelete}
	return userInfo, err
}

func (rr *redisRepository) RedisGetInviteInfo(token string) (string, error) {
	inviteInfoJson, err := rr.client.Get(token).Result()
	return inviteInfoJson, err
}

func (rr *redisRepository) RedisDelete(token string) {
	rr.client.Del(token)
}
