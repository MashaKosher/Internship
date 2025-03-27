package passwords

import "golang.org/x/crypto/bcrypt"

func HashPassword(passwrod string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ValidatePassword(HashedPass, RawPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(RawPass))
}
