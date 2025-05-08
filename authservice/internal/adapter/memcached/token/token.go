package token

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheTokenRepo struct {
	*memcache.Client
	logger di.LoggerType
}

func New(client *memcache.Client, logger di.LoggerType) *CacheTokenRepo {
	return &CacheTokenRepo{client, logger}
}

const tokenExpirationTime = 900 // time in secoinds

func (cr *CacheTokenRepo) AddTokenToCache(accessToken string, userId int) {

	key := getKey(accessToken)
	err := cr.Client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(strconv.Itoa(userId)),
		Expiration: tokenExpirationTime,
	})

	if err != nil {
		cr.logger.Error("Error while adding token to hash: " + err.Error())
	}
}

func (cr *CacheTokenRepo) GetTokenFromCache(accessToken string) (int, error) {

	key := getKey(accessToken)
	stringID, err := cr.Client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			cr.logger.Error("Token in cache not found")
			return -1, entity.ErrTokenInCacheNotFound
		}
		cr.logger.Error("Error while extracting token from cache: " + err.Error())
		return -1, err
	}

	id, err := strconv.Atoi(string(stringID.Value))
	if err != nil {
		return -1, err
	}

	return id, nil
}

// memcached has key constarint (<= 250 chars)
func getKey(accessToken string) string {
	if len(accessToken) < 250 {
		return accessToken
	}
	return accessToken[:250]

}
