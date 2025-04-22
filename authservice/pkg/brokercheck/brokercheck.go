package brokercheck

import (
	// repo "authservice/internal/adapter/db/sql"

	repo "authservice/internal/adapter/db/sql/auth"
	"authservice/pkg/logger"
	"authservice/pkg/tokens"
	"fmt"

	"authservice/internal/entity"

	producer "authservice/internal/adapter/kafka/producers"

	"github.com/golang-jwt/jwt"
)

func BrokerCheck(authRequest entity.AuthRequest) {
	var Answer entity.AuthAnswer

	validatedAccessToken, err := tokens.TokenVerify(authRequest.AccessToken)

	if err != nil {
		logger.Logger.Error("Inavlid " + tokens.ACCESS_TOKEN + " Token: " + err.Error())
		if err.Error() == "token expired" {
			validateRefreshToken(Answer, authRequest)
			return
		}

	}
	validateAccessToken(validatedAccessToken, Answer, authRequest)
}

func validateRefreshToken(Answer entity.AuthAnswer, authRequest entity.AuthRequest) {
	validatedRefreshToken, err := tokens.TokenVerify(authRequest.RefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedRefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	// Verifying Token type (it must be refresh)
	if err := tokens.VerifyTokenType(tokens.REFRESH_TOKEN, tokenType); err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	// Extracting User ID from valid refresh Token
	userId, err := tokens.GetIdFromValidatedToken(validatedRefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	logger.Logger.Info("User ID from " + tokens.REFRESH_TOKEN + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	DBUser, err := repo.FindUserById(int(userId))
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	logger.Logger.Info("User found successfully")

	// Creating Access Token
	newAccessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &DBUser)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	Answer.Role = string(DBUser.Role)
	Answer.ID = int32(DBUser.ID)
	Answer.Login = string(DBUser.Username)
	Answer.NewAccessToken = newAccessToken
	go producer.AnswerToken(Answer, authRequest.Partition)
	logger.Logger.Info("Access Token is expired, Refresg Token is valis, new Access Token" + fmt.Sprintln(Answer))
}

func validateAccessToken(validatedAccessToken *jwt.Token, Answer entity.AuthAnswer, authRequest entity.AuthRequest) {
	accessTokenType, err := tokens.GetTypeFromValidatedToken(validatedAccessToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	// Verifying Token type (it must be access)
	if err := tokens.VerifyTokenType(tokens.ACCESS_TOKEN, accessTokenType); err != nil {
		logger.Logger.Error(err.Error())
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	logger.Logger.Info("Token type is correct")

	userId, err := tokens.GetIdFromValidatedToken(validatedAccessToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	logger.Logger.Info("User ID from " + tokens.ACCESS_TOKEN + " Token is: " + fmt.Sprint(userId))

	DBUser, err := repo.FindUserById(int(userId))
	if err != nil {
		logger.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}
	logger.Logger.Info("User found successfully")

	Answer.Role = string(DBUser.Role)
	Answer.ID = int32(DBUser.ID)
	Answer.Login = string(DBUser.Username)
	logger.Logger.Info("Token is valid, " + fmt.Sprintln(Answer))
	go producer.AnswerToken(Answer, authRequest.Partition)
}

// func BrokerCheck(authRequest entity.AuthRequest) {

// 	var Answer entity.AuthAnswer

// 	validatedAccessToken, err := tokens.TokenVerify(authRequest.AccessToken)

// 	if err != nil {
// 		logger.Logger.Error("Inavlid " + tokens.ACCESS_TOKEN + " Token: " + err.Error())

// 		if err.Error() == "token expired" {
// 			validatedRefreshToken, err := tokens.TokenVerify(authRequest.RefreshToken)
// 			if err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}

// 			// Get Token type from Token
// 			tokenType, err := tokens.GetTypeFromValidatedToken(validatedRefreshToken)
// 			if err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}

// 			// Verifying Token type (it must be refresh)
// 			if err := tokens.VerifyTokenType(tokens.REFRESH_TOKEN, tokenType); err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}

// 			// Extracting User ID from valid refresh Token
// 			userId, err := tokens.GetIdFromValidatedToken(validatedRefreshToken)
// 			if err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}
// 			logger.Logger.Info("User ID from " + tokens.REFRESH_TOKEN + " Token is: " + fmt.Sprint(userId))

// 			// Search User by Id in DB
// 			DBUser, err := repo.FindUserById(int(userId))
// 			if err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}
// 			logger.Logger.Info("User found successfully")

// 			// Creating Access Token
// 			newAccessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &DBUser)
// 			if err != nil {
// 				Answer.Err = err.Error()
// 				go producer.AnswerToken(Answer, authRequest.Partition)
// 				return
// 			}
// 			Answer.Role = string(DBUser.Role)
// 			Answer.ID = int32(DBUser.ID)
// 			Answer.Login = string(DBUser.Username)
// 			Answer.NewAccessToken = newAccessToken
// 			go producer.AnswerToken(Answer, authRequest.Partition)
// 			logger.Logger.Info("Access Token is expired, Refresg Token is valis, new Access Token" + fmt.Sprintln(Answer))
// 			return

// 		}

// 	}

// 	accessTokenType, err := tokens.GetTypeFromValidatedToken(validatedAccessToken)
// 	if err != nil {
// 		logger.Logger.Error(err.Error())
// 		Answer.Err = err.Error()
// 		go producer.AnswerToken(Answer, authRequest.Partition)
// 		return
// 	}

// 	// Verifying Token type (it must be access)
// 	if err := tokens.VerifyTokenType(tokens.ACCESS_TOKEN, accessTokenType); err != nil {
// 		logger.Logger.Error(err.Error())
// 		Answer.Err = err.Error()
// 		go producer.AnswerToken(Answer, authRequest.Partition)
// 		return
// 	}
// 	logger.Logger.Info("Token type is correct")

// 	userId, err := tokens.GetIdFromValidatedToken(validatedAccessToken)
// 	if err != nil {
// 		logger.Logger.Error(err.Error())
// 		Answer.Err = err.Error()
// 		go producer.AnswerToken(Answer, authRequest.Partition)
// 		return
// 	}
// 	logger.Logger.Info("User ID from " + tokens.ACCESS_TOKEN + " Token is: " + fmt.Sprint(userId))

// 	DBUser, err := repo.FindUserById(int(userId))
// 	if err != nil {
// 		logger.Logger.Error("No User with such id: " + fmt.Sprint(userId))
// 		Answer.Err = err.Error()
// 		go producer.AnswerToken(Answer, authRequest.Partition)
// 		return
// 	}
// 	logger.Logger.Info("User found successfully")

// 	Answer.Role = string(DBUser.Role)
// 	Answer.ID = int32(DBUser.ID)
// 	Answer.Login = string(DBUser.Username)
// 	logger.Logger.Info("Token is valid, " + fmt.Sprintln(Answer))
// 	go producer.AnswerToken(Answer, authRequest.Partition)
// }
