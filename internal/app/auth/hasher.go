package auth

import (
	"golang.org/x/crypto/bcrypt"

	"DriverService/internal/app/appErrors"
)

type IHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(userPassword string, providedPassword string) error
}

type Hasher struct {
}

func (hasher *Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (hasher *Hasher) CheckPassword(userPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return appErrors.ErrWrongPassword
	}
	return nil
}
