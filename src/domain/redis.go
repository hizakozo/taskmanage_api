package domain

type RedisRepository interface {
	RedisSet(json string, key string)
	RedisGet(key string) (User, error)
	RedisGetInviteInfo(token string) (string, error)
	RedisDelete(token string)
}