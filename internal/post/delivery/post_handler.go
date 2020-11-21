package delivery

import (
	"net/http"
	"strconv"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/post"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	postUcase   post.PostUsecase
	userUcase   user.UserUsecase
	threadUcase thread.ThreadUsecase
	forumUcase  forum.ForumUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase,
	threadUcase thread.ThreadUsecase, forumUcase forum.ForumUsecase) *PostHandler {
	return &PostHandler{
		postUcase:   postUcase,
		userUcase:   userUcase,
		threadUcase: threadUcase,
		forumUcase:  forumUcase,
	}
}

func (ph *PostHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/post/:pid/details", ph.UpdatePostHandler())
}

func (ph *PostHandler) UpdatePostHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message" validate:"omitempty,gt=0"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		postID, _ := strconv.ParseUint(cntx.Param("pid"), 10, 64)
		postData := &models.Post{
			Message: req.Message,
		}

		post, err := ph.postUcase.Update(postID, postData)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, post)
	}
}
