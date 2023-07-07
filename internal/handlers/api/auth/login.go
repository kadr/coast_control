package auth

import (
	"context"
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/utils"
	"github.com/cost_control/internal/service/user"
	"github.com/cost_control/pkg/jwt"
	"github.com/cost_control/pkg/password_hasher"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IUserService interface {
	GetByEmail(ctx context.Context, email string) (user.UserServiceOutput, error)
}

type Auth struct {
	userService IUserService
	Response    utils.Response
	config      *config.Config
}

func New(userService IUserService, config *config.Config) *Auth {
	return &Auth{userService: userService, config: config}
}

func (a Auth) Login(c *gin.Context) {
	var request map[string]string
	if err := c.BindJSON(&request); err != nil {
		a.Response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var email string
	var ok bool
	if email, ok = request["email"]; !ok {
		a.Response.Error(c, http.StatusBadRequest, "не передан email")
		return
	}
	if _, ok = request["password"]; !ok {
		a.Response.Error(c, http.StatusBadRequest, "не передан пароль")
		return
	}
	findUser, err := a.userService.GetByEmail(context.Background(), email)
	if err != nil {
		a.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !password_hasher.Verify(findUser.Password, request["password"]) {
		a.Response.Error(c, http.StatusBadRequest, "не корректные email или пароль")
		return
	}
	token := jwt.New(findUser.Email, a.config.ExpiredAtMinutes)
	generatedToken, err := token.Generate(a.config.SignedKey)
	if err != nil {
		a.Response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	hour := 3600
	c.SetCookie("Authorization", generatedToken, hour, "", "", false, true)
	a.Response.Success(c, http.StatusOK, map[string]string{"token": generatedToken})
}
