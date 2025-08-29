package util

import (
	"testing"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	

	"github.com/golang-jwt/jwt/v5"
)


func TestGenerateAndVerifyJwt(t *testing.T) {

	JwtSecret= []byte("testsecret")

	userID := "user123"
	role := roles.Employee

	// Generate JWT
	tokenString, err := GenerateJwt(userID, role)
	if err != nil {
		t.Fatalf("GenerateJwt failed: %v", err)
	}
	if tokenString == "" {
		t.Fatal("GenerateJwt returned empty token")
	}

	// Verify JWT
	token, err := VarifyJwt(tokenString)
	if err != nil {
		t.Fatalf("VarifyJwt failed: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("failed to convert claims to MapClaims")
	}

	// Check the claims
	if claims["user_id"] != userID {
		t.Errorf("expected user_id %v, got %v", userID, claims["user_id"])
	}

	if int(claims["role"].(float64)) != int(role) {
		t.Errorf("expected role %v, got %v", role, claims["role"])
	}

	if claims["authorized"] != "true" {
		t.Errorf("expected authorized true, got %v", claims["authorized"])
	}

	// Check expiration is roughly 24h from now
	exp := int64(claims["exp"].(float64))
	if exp < time.Now().Unix()+23*3600 || exp > time.Now().Unix()+25*3600 {
		t.Errorf("unexpected expiration time: %v", exp)
	}

	// Test invalid token
	_, err =VarifyJwt("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}
