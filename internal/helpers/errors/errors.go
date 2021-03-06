package errors

import (
	"fmt"
	"net/http"

	. "github.com/OlegGibadulin/tech-db-forum/internal/consts"
)

type BodyType interface{}

type Error struct {
	Code     ErrorCode `json:"-"`
	HTTPCode int       `json:"-"`
	Body     BodyType  `json:"-"`
	Message  string    `json:"message"`
}

var WrongErrorCode = &Error{
	HTTPCode: http.StatusTeapot,
	Message:  "wrong error code",
}

func New(code ErrorCode, err error) *Error {
	customErr, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	copiedErr := *customErr
	copiedErr.Message = err.Error()
	return &copiedErr
}

func Get(code ErrorCode) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	return err
}

func BuildByBody(code ErrorCode, body BodyType) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	copiedErr := *err
	copiedErr.Body = body
	return &copiedErr
}

func BuildByMsg(code ErrorCode, attrs ...interface{}) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	copiedErr := *err
	copiedErr.Message = fmt.Sprintf(err.Message, attrs...)
	return &copiedErr
}

func (e *Error) Response() BodyType {
	if e.Message != "" {
		// return message responce
		return e
	}
	return e.Body
}

var Errors = map[ErrorCode]*Error{
	CodeBadRequest: {
		Code:     CodeBadRequest,
		HTTPCode: http.StatusNotFound,
		Message:  "Wrong request data",
	},
	CodeInternalError: {
		Code:     CodeInternalError,
		HTTPCode: http.StatusInternalServerError,
		Message:  "Something went wrong",
	},
	CodeUserAlreadyExists: {
		Code:     CodeUserAlreadyExists,
		HTTPCode: http.StatusConflict,
	},
	CodeUserDoesNotExist: {
		Code:     CodeUserDoesNotExist,
		HTTPCode: http.StatusNotFound,
		Message:  "Can't find user with %s %s",
	},
	CodeEmailAlreadyExists: {
		Code:     CodeEmailAlreadyExists,
		HTTPCode: http.StatusConflict,
		Message:  "User with email %s already exists",
	},
	CodeForumAlreadyExists: {
		Code:     CodeForumAlreadyExists,
		HTTPCode: http.StatusConflict,
	},
	CodeForumDoesNotExist: {
		Code:     CodeForumDoesNotExist,
		HTTPCode: http.StatusNotFound,
		Message:  "Can't find forum with %s %s",
	},
	CodeThreadAlreadyExists: {
		Code:     CodeThreadAlreadyExists,
		HTTPCode: http.StatusConflict,
	},
	CodeThreadDoesNotExist: {
		Code:     CodeThreadDoesNotExist,
		HTTPCode: http.StatusNotFound,
		Message:  "Can't find thread with %s %s",
	},
	CodeParentPostDoesNotExist: {
		Code:     CodeParentPostDoesNotExist,
		HTTPCode: http.StatusConflict,
		Message:  "Can't find parent post in thread with %s %d",
	},
	CodePostDoesNotExist: {
		Code:     CodePostDoesNotExist,
		HTTPCode: http.StatusNotFound,
		Message:  "Can't find post with %s %s",
	},
}
