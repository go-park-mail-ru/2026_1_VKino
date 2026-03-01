package auth

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_VKino/pkg/httpjson"
)

type httpErr struct {
	status  int
	message string
}

// Мапа "внутренняя ошибка -> внешний ответ"
var errToHTTP = map[error]httpErr{
	ErrUserAlreadyExists:  {status: http.StatusConflict, message: "user already exists"},
	ErrInvalidCredentials: {status: http.StatusUnauthorized, message: "invalid credentials"},
	ErrNoSession:          {status: http.StatusUnauthorized, message: "unauthorized"},
	ErrInvalidToken:       {status: http.StatusUnauthorized, message: "unauthorized"},
}

// перекладывать статус и тело в writeError
func writeServiceError(w http.ResponseWriter, err error) {
	var key error

	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		key = ErrUserAlreadyExists
	case errors.Is(err, ErrInvalidCredentials):
		key = ErrInvalidCredentials
	case errors.Is(err, ErrNoSession):
		key = ErrNoSession
	case errors.Is(err, ErrInvalidToken):
		key = ErrInvalidToken
	default:
		httpjson.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	mapped := errToHTTP[key]
	httpjson.WriteError(w, mapped.status, mapped.message)
}