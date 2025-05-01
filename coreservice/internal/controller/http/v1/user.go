package v1

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	u di.UserService
	l di.LoggerType
	v di.ValidatorType
}

func UserRoutes(router *gin.Engine, deps di.Container) {
	r := &userRoutes{u: deps.Services.User, l: deps.Logger, v: deps.Validator}
	user := router.Group("/user")
	user.POST("/deposit", r.makeDeposit)
	user.GET("/info", r.userInfo)
}

// MakeDeposit godoc
// @Summary Deposit funds to user account
// @Description Allows authenticated user to deposit funds to their balance
// @Tags User
// @Accept json
// @Produce json
// @Param deposit body entity.Balance true "Deposit amount details"
// @Success 200 {object} entity.Response "Returns new balance"
// @Failure 400 {object} entity.Error "Invalid token, negative amount or bad request"
// @Failure 401 {object} entity.Error "Unauthorized (missing or invalid token)"
// @Failure 500 {object} entity.Error "Internal server error"
// @Router /user/deposit [post]
func (r *userRoutes) makeDeposit(c *gin.Context) {
	data, exists := c.Get("data") // data contains DB user
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	var deposit entity.Balance
	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	err := r.v.Struct(deposit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	s := fmt.Sprint(deposit.Balance)
	parts := strings.Split(s, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "депозит должен иметь не более 2 знаков после запятой"})
		return
	}

	r.l.Info("Deposit: " + fmt.Sprintln(deposit))

	resp, err := r.u.MakeDeposit(data, deposit)

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
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
func (r *userRoutes) userInfo(c *gin.Context) {
	answer, exists := c.Get("data") // data contains DB user
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	user, err := r.u.UserInfo(answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}
