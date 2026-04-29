package authx

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
    return string(b), err
}

func VerifyPassword(hash, raw string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}
