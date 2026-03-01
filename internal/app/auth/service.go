package auth

import (
	"fmt"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userMap map[string]string
	userSessions map[string]TokenPair
}

const (
	// вынести в конфиг
	Secret = "4v29o7bg3p9h8cp9be7bc9w7bcg9py7bx9s"
	AccessTokenTTL = 15 * time.Minute
	RefreshTokenTTL = 14 * 24 * time.Hour
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoSession = errors.New("no session")
	ErrInvalidToken = errors.New("invalid token")
)

func NewService() *Service {
	// Саша, здесь нужно будет подключить твои модельки!
	// Мьютексы?
	return &Service{
		userMap: map[string]string{},
		userSessions: map[string]TokenPair{},
	}
}

// Саша, здесь логика на достать/создать, тоже меняй
func (s *Service) getUserByEmail(email string) (User, error) {
	passwordHash, exists := s.userMap[email]
	if !exists {
		return User{}, fmt.Errorf("no user")
	}
	return User{Email: email, Password: passwordHash}, nil
}

func (s *Service) createUser(email, passwordHash string) User {
	s.userMap[email] = passwordHash
	return User{Email: email, Password: passwordHash}
}
//

func (s *Service) userToString(user User) string {
	return user.Email
}

func (s *Service) tokenGenerate(user string, tokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject: user,
	})
	stringToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func (s *Service) tokenPairGenerate(user User) (TokenPair, error) {
	accessToken, err := s.tokenGenerate(s.userToString(user), AccessTokenTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("access token generate error: %w", err)
	}
	refreshToken, err := s.tokenGenerate(s.userToString(user), RefreshTokenTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("refresh token generate error: %w", err)
	}

	tokenPair := TokenPair{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	// переписать
	s.userSessions[user.Email] = tokenPair
	return tokenPair, nil
}


func (s *Service) SignIn(email, password string) (TokenPair, error) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}
	return s.tokenPairGenerate(user)
}


func (s *Service) SignUp(email, password string) (TokenPair, error) {
	_, err := s.getUserByEmail(email)
	if err == nil {
		return TokenPair{}, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, fmt.Errorf("bcrypt generate error: %w", err)
	}

	user := s.createUser(email, string(passwordHash))
	return s.tokenPairGenerate(user)
}


// обновляем протухший refresh-token
func (s *Service) refresh(email string) (TokenPair, error) {
	// переписать
	_, exists := s.userSessions[email]
	if !exists {
		return TokenPair{}, ErrNoSession
	}

	user, err := s.getUserByEmail(email)
	if err != nil {
		return TokenPair{}, ErrNoSession
	}

	return s.tokenPairGenerate(user)
}

func (s *Service) parseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *Service) validateAccessToken(tokenString string) (string, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims.Subject == "" {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}

func (s *Service) validateRefreshToken(tokenString string) (string, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims.Subject == "" {
		return "", ErrInvalidToken
	}

	// refresh должен совпадать с сохранённым у пользователя
	tokenPair, exists := s.userSessions[claims.Subject]
	if !exists {
		return "", ErrNoSession
	}

	if tokenPair.RefreshToken != tokenString {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}