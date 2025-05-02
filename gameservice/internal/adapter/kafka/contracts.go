package kafka

import (
	"gameservice/internal/entity"
)

type (
	AuthConsumer interface {
		Close()
		RecieveTokenInfo() (entity.AuthAnswer, error)
	}

	AuthProducer interface {
		Close()
		CheckAuthTokenRequest(accessToken, refreshToken string)
	}

	GameSettingsConsumer interface {
		Close()
		ConsumeGameSettings()
	}

	MatchInfoProducer interface {
		Close()
		SendMatchInfo(match entity.GameResult)
	}
)
