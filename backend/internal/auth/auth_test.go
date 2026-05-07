package auth

import (
	"strings"
	"testing"
)

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if strings.Contains(hash, "password123") {
		t.Fatal("hash must not contain plaintext password")
	}
	if !VerifyPassword("password123", hash) {
		t.Fatal("expected password to verify")
	}
	if VerifyPassword("wrong", hash) {
		t.Fatal("wrong password should not verify")
	}
}

func TestRegisterLoginAndAuthenticate(t *testing.T) {
	service := NewService()
	token, user, err := service.Register("user-1", "Test@Example.com", "password123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("email not normalized: %q", user.Email)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if _, ok := service.Sessions.sessions[token]; ok {
		t.Fatal("session repository must not store plaintext token as key")
	}
	loginToken, _, err := service.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if _, ok := service.Authenticate(loginToken); !ok {
		t.Fatal("expected login token to authenticate")
	}
}

func TestLoginRejectsBadPassword(t *testing.T) {
	service := NewService()
	if _, _, err := service.Register("user-1", "test@example.com", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, _, err := service.Login("test@example.com", "bad"); err == nil {
		t.Fatal("expected bad password to fail")
	}
}
