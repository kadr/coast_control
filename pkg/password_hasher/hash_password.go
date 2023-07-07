package password_hasher

import "golang.org/x/crypto/bcrypt"

func Gen(password string) (string, error) {
	return hash(password)
}

func Verify(dbPassword string, inputPassword string) bool {
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
