package core

import "golang.org/x/crypto/bcrypt"

type PasswordEncoder struct {
}

func (pe *PasswordEncoder) CompareHashAndPassword(hashedPassword string, password string) error {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
}

func (pe *PasswordEncoder) GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func NewPasswordEncoder() *PasswordEncoder {
	return &PasswordEncoder{}
}
