package auth

import (
	"fmt"
	"time"

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
		return TokenPair{}, fmt.Errorf("access token generate error")
	}
	refreshToken, err := s.tokenGenerate(s.userToString(user), RefreshTokenTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("refresh token generate error")
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
		return TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return TokenPair{}, fmt.Errorf("incorrect password")
	} else {
		return s.tokenPairGenerate(user)
	}
}


func (s *Service) SignUp(email, password string) (TokenPair, error) {
	_, err := s.getUserByEmail(email)
	if err == nil {
		return TokenPair{}, fmt.Errorf("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, fmt.Errorf("error in bcrypt")
	}

	user := s.createUser(email, string(passwordHash))
	return s.tokenPairGenerate(user)
}


// обновляем протухший refresh-token
func (s *Service) refresh(email string) (TokenPair, error) {
	// переписать
	_, exists := s.userSessions[email]
	if !exists {
		return TokenPair{}, fmt.Errorf("no session")
	}

	user, err := s.getUserByEmail(email)
	if err != nil {
		return TokenPair{}, err
	}

	return s.tokenPairGenerate(user)
}