package request_reader

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
	"github.com/OlegGibadulin/tech-db-forum/internal/helpers/errors"
	"github.com/OlegGibadulin/tech-db-forum/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RequestReader struct {
	cntx      echo.Context
	validator *validator.Validate
}

func NewRequestReader(cntx echo.Context) *RequestReader {
	return &RequestReader{
		cntx:      cntx,
		validator: validator.New(),
	}
}

func (rr *RequestReader) Read(request interface{}) *errors.Error {
	if err := rr.cntx.Bind(request); err != nil {
		return errors.New(CodeInternalError, err)
	}

	if err := rr.validator.Struct(request); err != nil {
		return errors.New(CodeBadRequest, err)
	}
	return nil
}

func (rr *RequestReader) ReadPosts() ([]*models.Post, *errors.Error) {
	b, err := ioutil.ReadAll(rr.cntx.Request().Body)
	if err != nil || b == nil {
		return nil, errors.New(CodeBadRequest, err)
	}

	posts := []*models.Post{}
	if err := json.Unmarshal(b, &posts); err != nil {
		return nil, errors.New(CodeBadRequest, err)
	}
	return posts, nil
}
