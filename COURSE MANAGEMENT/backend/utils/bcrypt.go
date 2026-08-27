package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mã hóa password
func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(bytes), err
}

// CheckPassword kiểm tra password
func CheckPassword(hash string, password string) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
