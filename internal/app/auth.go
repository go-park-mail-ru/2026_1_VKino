package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string
	Password string
}

type TokenPair struct {
    AccessToken  string
    RefreshToken string
}

const (
	Secret = "4v29o7bg3p9h8cp9be7bc9w7bcg9py7bx9s"
	AccessTokenTTL = 15 * time.Minute
	RefreshTokenTTL = 12 * time.Hour
)

var userMap = map[string]string {
	"user1": "password1",
	"user2": "password2",
	"user3": "password3",
} 

var userSessions = map[string]TokenPair {
	// "user1": TokenPair{},
}

func getUserByEmail(email string) (User, error) {
	passwordHash, exists := userMap[email]
	if !exists {
		return User{}, fmt.Errorf("no user")
	} 
	return User{Email: email, Password: passwordHash}, nil
}


func createUser(email, passwordHash string) User {
	userMap[email] = passwordHash
	return User{Email: email, Password: passwordHash}
}


func userToString(user User) string {
	return user.Email
}


func tokenGenerate(user string, tokenTTL time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject: user,
	})
	stringToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}


func tokenPairGenerate(user User) (TokenPair, error) {
	accessToken, err := tokenGenerate(userToString(user), AccessTokenTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("access token generate error")
	}
	refreshToken, err := tokenGenerate(userToString(user), RefreshTokenTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("refresh token generate error")
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}


func auth(email, password string) (TokenPair, error) {
	// go to db and get User
	user, err := getUserByEmail(email)
	if err != nil {
		return TokenPair{}, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, fmt.Errorf("error in bcrypt")
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(user.Password))
	if err != nil {
		return TokenPair{}, fmt.Errorf("incorrect password")
	} else {
		// выдаём access + refresh jwt tokens
		return tokenPairGenerate(user)
	}

}


func register(email, password string) (TokenPair, error) {
	_, err := getUserByEmail(email)
	if err == nil {
		return TokenPair{}, fmt.Errorf("user already exists")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, fmt.Errorf("error in bcrypt")
	}
	// go to db and create user
	user := createUser(email, string(passwordHash))
	return tokenPairGenerate(user)
}


func middleware() {

}


func main() {

	//auth test
	user_email := "user1"
	user_password := "password2"
	
	token, err := auth(user_email, user_password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	//register test
	user_email = "user4"
	user_password = "password4"
	token, err = register(user_email, user_password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

}