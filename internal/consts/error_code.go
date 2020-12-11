package consts

type ErrorCode uint16

const (
	CodeBadRequest ErrorCode = iota + 101
	CodeInternalError
	CodeUserAlreadyExists
	CodeUserDoesNotExist
	CodeEmailAlreadyExists
	CodeForumAlreadyExists
	CodeForumDoesNotExist
	CodeThreadAlreadyExists
	CodeThreadDoesNotExist
	CodeParentPostDoesNotExist
	CodePostDoesNotExist
)

const OnPostInsertExceptionMsg = "pq: Can not find parent post into thread"
