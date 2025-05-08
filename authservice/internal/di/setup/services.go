package setup

import (
	authRepo "authservice/internal/adapter/db/sql/auth"
	kafkaRepo "authservice/internal/adapter/kafka"
	"authservice/internal/di"
	"authservice/internal/usecase/auth"
)

func mustServices(db di.DBType, logger di.LoggerType, RSAKeys di.RSAKeys, signUpProducer kafkaRepo.SignUpProducer, cache di.Cache) di.Services {

	authUseCase := createAuthUseCase(db, logger, RSAKeys, signUpProducer, cache)

	return di.Services{
		Auth: authUseCase,
	}
}

func createAuthUseCase(db di.DBType, logger di.LoggerType, RSAKeys di.RSAKeys, signUpProducer kafkaRepo.SignUpProducer, cache di.Cache) *auth.UseCase {
	return auth.New(
		authRepo.New(db),
		logger,
		RSAKeys,
		signUpProducer,
		cache,
	)
}
