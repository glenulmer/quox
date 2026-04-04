package globals

import (
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
    if password == "" { return "" }
    bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes)
}

func Authenticate(dbPass, userIn string) bool {
    e := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(userIn))
    return e == nil
}