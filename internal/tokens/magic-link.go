package tokens

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"time"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
)

func CreateMagicLinkToken(email string) (string, error) {
	var value = map[string]string{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 60).Format(time.RFC3339),
	}

	return cookieHandler.Encode("magic-link", value)
}

func ValidateMagicLinkToken(token string) (string, error) {
	var value = make(map[string]string)

	err := cookieHandler.Decode("magic-link", token, &value)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	if value["exp"] == "" {
		return "", fmt.Errorf("invalid token")
	}

	exp, err := time.Parse(time.RFC3339, value["exp"])
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	if time.Now().After(exp) {
		return "", fmt.Errorf("token expired")
	}

	return value["email"], nil
}