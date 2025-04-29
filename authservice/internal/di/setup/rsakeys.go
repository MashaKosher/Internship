package setup

import (
	"authservice/internal/di"
	"authservice/pkg/keys"
)

func mustRSAKeys(cfg di.ConfigType, logger di.LoggerType) di.RSAKeys {
	var RSAKeys di.RSAKeys
	keys.ReadRSAKeys(cfg, logger, &RSAKeys)
	return RSAKeys
}
