package password_hasher

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct{}

func New() *PasswordHasher {
	return &PasswordHasher{}
}

func (ph PasswordHasher) Gen(password string) (string, error) {
	return hash(password)
}

func (ph PasswordHasher) Verify(dbPassword string, inputPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword)); err != nil {
		return false
	}

	return true
}

func hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
