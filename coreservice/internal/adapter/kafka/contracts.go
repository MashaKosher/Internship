package kafka

import (
	"coreservice/internal/entity"
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

	DailyTaskConsumer interface {
		Close()
		ReceiveDailyTask()
	}

	MatchInfoConsumer interface {
		Close()
		RecieveMatchInfo()
	}

	SeasonInfoConsumer interface {
		Close()
		RecieveSeasonInfo()
	}

	UserSignUpConsumer interface {
		Close()
		ReceiveSignedUpUser()
	}
)
