package delivery

import (
	"net/http"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUcase user.UserUsecase
}

func NewUserHandler(userUcase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUcase: userUcase,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/user/:nickname/create", uh.CreateUserHandler())
	e.GET("/api/user/:nickname/profile", uh.GetUserHandler())
	e.POST("/api/user/:nickname/profile", uh.UpdateUserHandler())
}

func (uh *UserHandler) CreateUserHandler() echo.HandlerFunc {
	type Request struct {
		models.User
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if err := uh.userUcase.Create(&req.User); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, req.User)
	}
}

func (uh *UserHandler) UpdateUserHandler() echo.HandlerFunc {
	type Request struct {
		Fullname string `json:"fullname" validate:"omitempty,gte=3,lte=32"`
		Email    string `json:"email" validate:"omitempty,email,lte=64"`
		About    string `json:"about"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		nickname := cntx.Param("nickname")
		userData := &models.User{
			Fullname: req.Fullname,
			Email:    req.Email,
			About:    req.About,
		}

		user, err := uh.userUcase.Update(nickname, userData)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, user)
	}
}

func (uh *UserHandler) GetUserHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		nickname := cntx.Param("nickname")
		user, err := uh.userUcase.GetByNickname(nickname)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, user)
	}
}
