package delivery

import (
	"net/http"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ThreadHandler struct {
	threadUcase thread.ThreadUsecase
	userUcase   user.UserUsecase
	forumUcase  forum.ForumUsecase
}

func NewThreadHandler(threadUcase thread.ThreadUsecase, userUcase user.UserUsecase, forumUcase forum.ForumUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUcase: threadUcase,
		userUcase:   userUcase,
		forumUcase:  forumUcase,
	}
}

func (th *ThreadHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/forum/:forum/create", th.CreateThreadHandler())
}

func (th *ThreadHandler) CreateThreadHandler() echo.HandlerFunc {
	type Request struct {
		models.Thread
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if _, err := th.forumUcase.GetBySlug(req.Forum); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if _, err := th.userUcase.GetByNickname(req.Author); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if err := th.threadUcase.Create(&req.Thread); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, req.Thread)
	}
}
