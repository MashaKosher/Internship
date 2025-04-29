package setup

import (
	authRepo "authservice/internal/adapter/db/sql/auth"
	"authservice/internal/di"
	"authservice/internal/usecase/auth"
)

func mustServices(db di.DBType, logger di.LoggerType, RSAKeys di.RSAKeys) di.Services {

	authUseCase := auth.New(
		authRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
		logger,
		RSAKeys,
	)

	return di.Services{
		Auth: authUseCase,
	}
}
