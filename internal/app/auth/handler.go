package auth

import (
	"net/http"
	"errors"

	"github.com/go-park-mail-ru/2026_1_VKino/pkg/httpjson"
)

type Handler struct {
	service *Service
}


func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}


func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /sign-up", h.signUp)
	mux.HandleFunc("POST /sign-in", h.signIn)
	mux.HandleFunc("POST /refresh", h.refresh)
}


func writeError(w http.ResponseWriter, status int, message string) {
	httpjson.WriteJSON(w, status, errorResponse{Error: message})
}


// обработка ошибок от service.go
func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		writeError(w, http.StatusConflict, "user already exists")
	case errors.Is(err, ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, ErrNoSession), errors.Is(err, ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}


// ставим refresh-token в cookie
func setRefreshCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(RefreshTokenTTL.Seconds()),
	})
}


func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := httpjson.ReadJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	tokens, err := h.service.SignUp(req.Email, req.Password)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	setRefreshCookie(w, tokens.RefreshToken)
	httpjson.WriteJSON(w, http.StatusCreated, tokens)
}


func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := httpjson.ReadJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	tokens, err := h.service.SignIn(req.Email, req.Password)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	setRefreshCookie(w, tokens.RefreshToken)
	httpjson.WriteJSON(w, http.StatusOK, tokens)
}


func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	email, err := h.service.validateRefreshToken(cookie.Value)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	tokenPair, err := h.service.refresh(email)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	setRefreshCookie(w, tokenPair.RefreshToken)
	httpjson.WriteJSON(w, http.StatusOK, accessTokenResponse{
		AccessToken: tokenPair.AccessToken,
	})
}