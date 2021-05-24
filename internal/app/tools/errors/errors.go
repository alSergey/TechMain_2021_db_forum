package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorType uint8

const (
	WrongErrorCode ErrorType = iota + 1
	InternalError
	UserCreateExist
	UserProfileNotExist
	UserProfileConflict
	ForumCreateNotExist
	ForumCreateConflict
	ForumDetailsNotExist
	ForumCreateThreadNotExist
	ForumCreateThreadConflict
	ForumThreadsNotExist
)

type Error struct {
	ErrorCode ErrorType `json:"-"`
	HttpError int       `json:"-"`
	Message   string    `json:"message"`
}

func JSONError(error *Error, w http.ResponseWriter) {
	body, err := json.Marshal(error)
	if err != nil {
		return
	}

	w.WriteHeader(error.HttpError)
	w.Write(body)
}

func JSONSuccess(status int, data interface{}, w http.ResponseWriter) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}

var CustomErrors = map[ErrorType]*Error{
	WrongErrorCode: {
		ErrorCode: WrongErrorCode,
		HttpError: http.StatusInternalServerError,
	},
	InternalError: {
		ErrorCode: InternalError,
		HttpError: http.StatusInternalServerError,
		Message:   "something wrong",
	},
	UserCreateExist: {
		ErrorCode: UserCreateExist,
		HttpError: http.StatusConflict,
	},
	UserProfileNotExist: {
		ErrorCode: UserProfileNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Can't find user\n",
	},
	UserProfileConflict: {
		ErrorCode: UserProfileConflict,
		HttpError: http.StatusConflict,
		Message:   "Conflict email\n",
	},
	ForumCreateNotExist: {
		ErrorCode: ForumCreateNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Can't find user\n",
	},
	ForumCreateConflict: {
		ErrorCode: ForumCreateConflict,
		HttpError: http.StatusConflict,
	},
	ForumDetailsNotExist: {
		ErrorCode: ForumDetailsNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Can't find forum\n",
	},
	ForumCreateThreadNotExist: {
		ErrorCode: ForumCreateThreadNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Can't find user\n",
	},
	ForumCreateThreadConflict: {
		ErrorCode: ForumCreateThreadConflict,
		HttpError: http.StatusConflict,
	},
	ForumThreadsNotExist: {
		ErrorCode: ForumThreadsNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Can't find forum\n",
	},
}

func Cause(code ErrorType) *Error {
	err, ok := CustomErrors[code]
	if !ok {
		return CustomErrors[WrongErrorCode]
	}

	return err
}

func UnexpectedInternal(err error) *Error {
	unexpErr := CustomErrors[InternalError]
	unexpErr.Message = err.Error()

	return unexpErr
}
