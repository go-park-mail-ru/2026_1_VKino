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

// Порядок важен, чтобы поведение было детерминированным
var errOrder = []error{
	ErrUserAlreadyExists,
	ErrInvalidCredentials,
	ErrNoSession,
	ErrInvalidToken,
}

// перекладывать статус и тело в writeError
func writeServiceError(w http.ResponseWriter, err error) {
	for _, target := range errOrder {
		if errors.Is(err, target) {
			mapped := errToHTTP[target]
			httpjson.WriteError(w, mapped.status, mapped.message)
			return
		}
	}

	httpjson.WriteError(w, http.StatusInternalServerError, "internal server error")
}