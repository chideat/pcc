package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	JWT_AUTH_KEY = "5233c3ee93e722290f6de503d64665c9835f30b483107e93dec31c8ff7dadcc1"
)

func Rand() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}

func PasswordEncrypt(password string) (string, string) {
	salt := Rand()
	ePwd := sha256.Sum256([]byte(fmt.Sprintf("%s-%s", password, salt)))
	return fmt.Sprintf("%x", ePwd), salt
}

func PasswordValid(password, salt, ePwd string) bool {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s-%s", password, salt)))) == ePwd
}

type ThirdAuth struct {
	Id     string
	Secret string
}

type UserTokenInfo struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Expired time.Time `json:"expired"`
}

func (info *UserTokenInfo) String() string {
	bytes, _ := json.Marshal(info)
	return fmt.Sprintf("%s", bytes)
}

func Auth(authToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(JWT_AUTH_KEY), nil
		/*
			get user by user id
			userID, ok := token.Claims["id"].(int64)
			if !ok {
				return nil, fmt.Errorf("invalid user id")
			}
			user, err := GetUserById(userID)
			if err != nil {
				return nil, fmt.Errorf("invalid user id")
			}
			token["user"] = user

			return []byte(user.SigningKey), nil
		*/
	})
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}
