package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name            string
		header          *string
		wantTokenString string
		wantErr         bool
	}{
		{
			name:            "Valid header",
			header:          new("Bearer 123456"),
			wantTokenString: "123456",
			wantErr:         false,
		},
		{
			name:            "Valid header with whitespaces",
			header:          new("  Bearer  123456 "),
			wantTokenString: "123456",
			wantErr:         false,
		},
		{
			name:            "No whitespace between prefix and token",
			header:          new("  Bearer123456 "),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "Invalid header prefix",
			header:          new("Beerer 123456"),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "Missing prefix",
			header:          new("123456"),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "Missing token",
			header:          new("Bearer "),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "Empty header",
			header:          new(""),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "Missing header",
			header:          nil,
			wantTokenString: "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpHeader := http.Header{}
			if tt.header != nil {
				httpHeader.Add("Authorization", *tt.header)
			}

			gotToken, err := GetBearerToken(httpHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantTokenString {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantTokenString)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token := MakeResfreshToken()
	if len(token) != 64 {
		t.Errorf("Expected token to be 64 characters long, gotToken = %v, gotLength = %v", token, len(token))
	}
}
