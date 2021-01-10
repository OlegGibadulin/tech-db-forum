package delivery

import (
	"net/http"
	"time"

	"github.com/OlegGibadulin/tech-db-forum/internal/forum"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
)

type ForumHandler struct {
	forumUcase  forum.ForumUsecase
	userUcase   user.UserUsecase
	threadUcase thread.ThreadUsecase
}

func NewForumHandler(forumUcase forum.ForumUsecase, userUcase user.UserUsecase,
	threadUcase thread.ThreadUsecase) *ForumHandler {
	return &ForumHandler{
		forumUcase:  forumUcase,
		userUcase:   userUcase,
		threadUcase: threadUcase,
	}
}

func (fh *ForumHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/forum/create", fh.CreateForumHandler())
	e.GET("/api/forum/:slug/details", fh.GetForumDetailesHandler())
	e.POST("/api/forum/:forum/create", fh.CreateThreadHandler())
	e.GET("/api/forum/:slug/threads", fh.GetThreadsByForumHandler())
	e.GET("/api/forum/:slug/users", fh.GetUsersByForumHandler())
}

func (fh *ForumHandler) CreateForumHandler() echo.HandlerFunc {
	type Request struct {
		models.Forum
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		user, err := fh.userUcase.GetByNickname(req.User)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		req.User = user.Nickname

		if err := fh.forumUcase.Create(&req.Forum); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, req.Forum)
	}
}

func (fh *ForumHandler) GetForumDetailesHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		slug := cntx.Param("slug")
		forum, err := fh.forumUcase.GetBySlug(slug)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, forum)
	}
}

func (fh *ForumHandler) CreateThreadHandler() echo.HandlerFunc {
	type Request struct {
		models.Thread
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		forum, err := fh.forumUcase.GetBySlug(req.Forum)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		req.Forum = forum.Slug

		user, err := fh.userUcase.GetByNickname(req.Author)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		req.Author = user.Nickname

		if err := fh.threadUcase.Create(&req.Thread); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, req.Thread)
	}
}

func (fh *ForumHandler) GetThreadsByForumHandler() echo.HandlerFunc {
	type Request struct {
		Since time.Time `query:"since"`
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slug := cntx.Param("slug")
		if _, err := fh.forumUcase.GetBySlug(slug); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		threads, err := fh.threadUcase.ListByForum(slug, req.Since, &req.Pagination)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, threads)
	}
}

func (fh *ForumHandler) GetUsersByForumHandler() echo.HandlerFunc {
	type Request struct {
		Since string `query:"since"`
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slug := cntx.Param("slug")
		if _, err := fh.forumUcase.GetBySlug(slug); err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		users, err := fh.userUcase.ListByForum(slug, req.Since, &req.Pagination)
		if err != nil {
			// logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, users)
	}
}
