package service

import (
	"fmt"
	"math"
	"net/http"

	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"coreservice/pkg"

	"coreservice/internal/repository/sqlc"

	"github.com/gin-gonic/gin"
)

// CheckToken godoc
// @Summary Validate JWT token
// @Description Verifies JWT token validity and returns user information
// @Tags Authentication
// @Produce json
// @Success 200 {object} entity.TypeResponse "Token validation response with user data"
// @Failure 400 {object} entity.Error "Invalid token format"
// @Failure 401 {object} entity.Error "Missing or invalid token"
// @Router /check-token [get]
func CheckToken(c *gin.Context) {
	message, exists := c.Get("message")
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
		return
	}

	data, exists := c.Get("data")
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	user, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	var resp entity.TypeResponse

	resp.User = pkg.GetUserInfo(&user)
	resp.Message = message.(string)

	c.JSON(http.StatusOK, resp)

}

// MakeDeposit godoc
// @Summary Deposit funds to user account
// @Description Allows authenticated user to deposit funds to their balance
// @Tags User
// @Accept json
// @Produce json
// @Param deposit body entity.BalanceBody true "Deposit amount details"
// @Success 200 {object} entity.Response "Returns new balance"
// @Failure 400 {object} entity.Error "Invalid token, negative amount or bad request"
// @Failure 401 {object} entity.Error "Unauthorized (missing or invalid token)"
// @Failure 500 {object} entity.Error "Internal server error"
// @Router /user/deposit [post]
func MakeDeposit(c *gin.Context) {
	data, exists := c.Get("data") // data contains DB user
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	var deposit entity.BalanceBody
	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}
	logger.Logger.Info("Deposit: " + fmt.Sprintln(deposit))

	if deposit.Balance < 0 {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Deposit cannot be less than zero"})
		return
	}

	player, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	var currentBalance float64
	if player.Balance.Valid { // Chcecking if balance is NOT NULL
		left, _ := player.Balance.Int.Float64()
		currentBalance = left / math.Pow(float64(10), -float64(player.Balance.Exp))
	}

	player, err = sqlc.UpdateBalance(player.ID, deposit.Balance+currentBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity.Response{Message: player.Balance.Int.String()})
}

// UserInfo godoc
// @Summary Get authenticated user information
// @Description Returns current user's information based on valid JWT token
// @Tags User
// @Produce json
// @Success 200 {object} entity.User "User information"
// @Failure 400 {object} entity.Error "Invalid token or conversion error"
// @Failure 401 {object} entity.Error "Unauthorized (when token is missing)"
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	answer, exists := c.Get("data") // data contains DB user
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	user, err := pkg.ConvertAnyToDBUser(answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pkg.GetUserInfo(&user))

}
