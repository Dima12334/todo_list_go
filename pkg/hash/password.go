package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	GeneratePasswordHash(password string) (string, error)
	CheckPasswordHash(hash, password string) (bool, error)
}

type SHA1Hasher struct{}

func NewSHA1Hasher() *SHA1Hasher {
	return &SHA1Hasher{}
}

func (h *SHA1Hasher) GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *SHA1Hasher) CheckPasswordHash(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}
