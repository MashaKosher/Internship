package setup

import (
	authRepo "authservice/internal/adapter/db/sql/auth"
	kafkaRepo "authservice/internal/adapter/kafka"
	"authservice/internal/di"
	"authservice/internal/usecase/auth"
)

func mustServices(db di.DBType, logger di.LoggerType, RSAKeys di.RSAKeys, signUpProducer kafkaRepo.SignUpProducer) di.Services {

	authUseCase := createAuthUseCase(db, logger, RSAKeys, signUpProducer)

	return di.Services{
		Auth: authUseCase,
	}
}

func createAuthUseCase(db di.DBType, logger di.LoggerType, RSAKeys di.RSAKeys, signUpProducer kafkaRepo.SignUpProducer) *auth.UseCase {
	return auth.New(
		authRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
		logger,
		RSAKeys,
		signUpProducer,
	)
}
