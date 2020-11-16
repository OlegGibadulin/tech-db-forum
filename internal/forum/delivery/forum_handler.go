package delivery

import (
	"net/http"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ForumHandler struct {
	forumUcase forum.ForumUsecase
	userUcase  user.UserUsecase
}

func NewForumHandler(forumUcase forum.ForumUsecase, userUcase user.UserUsecase) *ForumHandler {
	return &ForumHandler{
		forumUcase: forumUcase,
		userUcase:  userUcase,
	}
}

func (uh *ForumHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/forum/create", uh.CreateForumHandler())
}

func (uh *ForumHandler) CreateForumHandler() echo.HandlerFunc {
	type Request struct {
		models.Forum
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if _, err := uh.userUcase.GetByNickname(req.User); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if err := uh.forumUcase.Create(&req.Forum); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, req.Forum)
	}
}
