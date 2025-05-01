package middleware

import (
	"coreservice/internal/adapter/kafka/consumers"
	"coreservice/internal/adapter/kafka/producers"
	"coreservice/internal/di"
	"coreservice/pkg"

	"fmt"
	"net/http"

	userRepo "coreservice/internal/adapter/db/postgres/user"
	userNameElasticRepo "coreservice/internal/adapter/elastic/user"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(cfg di.ConfigType, logger di.LoggerType, db di.DBType, ESClient di.ESClient, Index di.ElasticIndex, bus di.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {

		userRepo := userRepo.New(db)
		elastic := userNameElasticRepo.New(ESClient, Index, logger, userRepo)

		logger.Info("Middleware for Token check")
		accessToken, err := c.Cookie("access")
		if err != nil || len(accessToken) == 0 {
			logger.Error("No access token")
			c.JSON(http.StatusBadRequest, gin.H{"error": "No access token"})
			c.Abort()
			return
		}
		logger.Info("Access Token:" + accessToken)

		refreshToken, _ := c.Cookie("refresh")

		// Sending token to authserver through kafka
		producers.CheckToken(accessToken, refreshToken, cfg, bus)

		// Receiving an answer from authserver through kafka
		answer, _ := consumers.RecieveTokenInfo(cfg, bus)

		// If token is not valid fieвs Err will be not empty
		if len(answer.Err) != 0 {
			logger.Error(answer.Err)
			c.JSON(http.StatusBadRequest, gin.H{"error": answer.Err})
			c.Abort()
			return
		}

		// If refresh token is valid, we recieve new access
		if len(answer.NewAccessToken) != 0 {
			logger.Info("New access Token")
			c.SetCookie("access", answer.NewAccessToken, 3600, "/", "", false, true)
		}

		// Serching for user in DB using Elastic
		player, exists := userRepo.GetPlayerById(answer.ID)
		if exists {
			logger.Info("Player already in DB: " + fmt.Sprintln(player))
			c.Set("message", "User already in DB")
			c.Set("data", player)
			c.Next()
			return
		}

		// If user not exists adding user to DB
		player, err = userRepo.AddPlayer(answer)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Adding new user to Elastic Index
		user := pkg.GetUserInfo(&player)

		if err := elastic.AddUserToIndex(user, int(player.ID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		logger.Info("User added successfully: " + fmt.Sprintln(player))
		c.Set("message", "Player added succesfully")
		c.Set("data", player)
		c.Next()
	}
}
