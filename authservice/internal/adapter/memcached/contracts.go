package memcached

type (
	TokenCacheRepo interface {
		AddTokenToCache(accessToken string, userId int)
		GetTokenFromCache(accessToken string) (int, error)
	}
)
