package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/post"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
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
	e.GET("/api/post/:pid/details", ph.GetPostDetailesHandler())
	e.POST("/api/post/:pid/details", ph.UpdatePostHandler())
}

func (ph *PostHandler) UpdatePostHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message" validate:"omitempty,gt=0"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		postID, _ := strconv.ParseUint(cntx.Param("pid"), 10, 64)
		postData := &models.Post{
			Message: req.Message,
		}

		post, err := ph.postUcase.Update(postID, postData)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, post)
	}
}

func (ph *PostHandler) GetPostDetailesHandler() echo.HandlerFunc {
	type Response struct {
		Post   *models.Post   `json:"post"`
		Author *models.User   `json:"author"`
		Thread *models.Thread `json:"thread"`
		Forum  *models.Forum  `json:"forum"`
	}

	return func(cntx echo.Context) error {
		related := strings.Split(cntx.QueryParam("related"), ",")
		postID, _ := strconv.ParseUint(cntx.Param("pid"), 10, 64)

		res := &Response{}
		var err *errors.Error

		if res.Post, err = ph.postUcase.GetByID(postID); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		for _, param := range related {
			switch param {
			case "user":
				if res.Author, err = ph.userUcase.GetByPostID(postID); err != nil {
					// logrus.Error(err.Message)
					return cntx.JSON(err.HTTPCode, err.Response())
				}
			case "forum":
				if res.Forum, err = ph.forumUcase.GetByPostID(postID); err != nil {
					// logrus.Error(err.Message)
					return cntx.JSON(err.HTTPCode, err.Response())
				}
			case "thread":
				if res.Thread, err = ph.threadUcase.GetByPostID(postID); err != nil {
					// logrus.Error(err.Message)
					return cntx.JSON(err.HTTPCode, err.Response())
				}
			}
		}
		return cntx.JSON(http.StatusOK, res)
	}
}
