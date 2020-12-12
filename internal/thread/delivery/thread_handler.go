package delivery

import (
	"net/http"

	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
	"github.com/OlegGibadulin/tech-db-forum/internal/post"
	"github.com/OlegGibadulin/tech-db-forum/internal/thread"
	"github.com/OlegGibadulin/tech-db-forum/internal/user"
	"github.com/OlegGibadulin/tech-db-forum/pkg/uniq"
	reader "github.com/OlegGibadulin/tech-db-forum/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ThreadHandler struct {
	threadUcase thread.ThreadUsecase
	userUcase   user.UserUsecase
	postUcase   post.PostUsecase
}

func NewThreadHandler(threadUcase thread.ThreadUsecase, userUcase user.UserUsecase,
	postUcase post.PostUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUcase: threadUcase,
		userUcase:   userUcase,
		postUcase:   postUcase,
	}
}

func (th *ThreadHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.GET("/api/thread/:slug_or_id/details", th.GetThreadDetailesHandler())
	e.POST("/api/thread/:slug_or_id/details", th.UpdateThreadHandler())
	e.POST("/api/thread/:slug_or_id/vote", th.VoteThreadHandler())
	e.POST("/api/thread/:slug_or_id/create", th.CreatePostsHandler())
	e.GET("/api/thread/:slug_or_id/posts", th.GetPostsByThreadHandler())
}

func (th *ThreadHandler) UpdateThreadHandler() echo.HandlerFunc {
	type Request struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slugOrID := cntx.Param("slug_or_id")
		threadData := &models.Thread{
			Title:   req.Title,
			Message: req.Message,
		}

		thread, err := th.threadUcase.Update(slugOrID, threadData)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) GetThreadDetailesHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		slugOrID := cntx.Param("slug_or_id")
		thread, err := th.threadUcase.GetBySlugOrID(slugOrID)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) VoteThreadHandler() echo.HandlerFunc {
	type Request struct {
		models.Vote
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if _, err := th.userUcase.GetByNickname(req.Nickname); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slugOrID := cntx.Param("slug_or_id")
		thread, err := th.threadUcase.Vote(slugOrID, &req.Vote)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) CreatePostsHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		posts, err := reader.NewRequestReader(cntx).ReadPosts()
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slugOrID := cntx.Param("slug_or_id")
		thread, err := th.threadUcase.GetBySlugOrID(slugOrID)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		var nicknames []string
		for _, post := range posts {
			nicknames = append(nicknames, post.Author)
		}
		nicknames = uniq.RemoveDuplicates(nicknames)
		if err := th.userUcase.CheckUsersExistence(nicknames); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		if err := th.postUcase.Create(posts, thread); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusCreated, posts)
	}
}

func (th *ThreadHandler) GetPostsByThreadHandler() echo.HandlerFunc {
	type Request struct {
		Since uint64 `query:"since"`
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		slugOrID := cntx.Param("slug_or_id")
		threadID, err := th.threadUcase.CheckThreadExistence(slugOrID)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}

		posts, err := th.postUcase.ListByThread(threadID, req.Since, &req.Pagination)
		if err != nil {
			logrus.Error(err.Message)
			return cntx.JSON(err.HTTPCode, err.Response())
		}
		return cntx.JSON(http.StatusOK, posts)
	}
}
