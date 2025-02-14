package auth

import (
	"net/http"
	"testing"
)

func TestHashPassword(t *testing.T) {

	password := "Salam45"

	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword failed :%v", err)
	}

	if hash == "" {
		t.Errorf("HashPassword returned empty hash")
	}

	if hash == password {
		t.Errorf("Fail, hash equal passordd")
	}

}

func TestCheckPasswordHash(t *testing.T) {

	password := "Salam45"

	hash, _ := HashPassword(password)

	err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Different hash and passord :%v", err)
	}

}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{
		"Authorization": {"Bearer TOKEN_STRING"},
	}

	testString, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no errror, bit got: %v", err)
	}

	expectedToken := "TOKEN_STRING"
	if testString != expectedToken {
		t.Errorf("expected token to be %s, but got %s", expectedToken, testString)
	}
}
