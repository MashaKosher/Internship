package kafka

import "adminservice/internal/entity"

type (
	AuthConsumer interface {
		Close()
		AnswerTokens() (entity.AuthAnswer, error)
	}

	AuthProducer interface {
		Close()
		CheckToken(accessToken, refreshToken string)
	}

	GameSettingsProducer interface {
		Close()
		SendGameSettings(season entity.SettingsJson)
	}

	SeasonProducer interface {
		Close()
		SendSeasonInfo(season entity.SeasonOut)
	}

	DailyTaskProducer interface {
		Close()
		SendDailyTask(task entity.DailyTasks)
	}
)
