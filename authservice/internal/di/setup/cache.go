package setup

import (
	"authservice/internal/adapter/memcached/token"
	"authservice/internal/di"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func mustCache(cfg di.ConfigType, logger di.LoggerType) di.Cache {

	client := memcache.New(fmt.Sprintf("%s:%s", cfg.Memcached.Host, cfg.Memcached.Port))

	return di.Cache{
		Token: token.New(client, logger),
	}
}
