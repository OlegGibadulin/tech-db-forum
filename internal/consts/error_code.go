package consts

type ErrorCode uint16

const (
	CodeBadRequest ErrorCode = iota + 101
	CodeInternalError
	CodeUserAlreadyExists
	CodeUserDoesNotExist
	CodeEmailAlreadyExists
)
