package auth_test

import (
	"ValenTheRed/chirpy/internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
)

const defaultHmacKey = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func TestValidateJWT(t *testing.T) {
	uuids := uuid.UUIDs{
		uuid.New(),
		uuid.New(),
	}
	makeJWTArgs := []struct {
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		{
			userID:      uuids[0],
			tokenSecret: defaultHmacKey,
			expiresIn:   time.Hour,
		},
		{
			userID:      uuids[1],
			tokenSecret: defaultHmacKey,
			expiresIn:   -time.Hour,
		},
	}
	tokens := make([]string, len(makeJWTArgs))
	for i, args := range makeJWTArgs {
		token, _ := auth.MakeJWT(args.userID, args.tokenSecret, args.expiresIn)
		tokens[i] = token
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "token is valid",
			tokenString: tokens[0],
			tokenSecret: defaultHmacKey,
			want:        uuids[0],
			wantErr:     false,
		},
		{
			name:        "token is invalid",
			tokenString: "",
			tokenSecret: defaultHmacKey,
			want:        uuids[0],
			wantErr:     true,
		},
		{
			name:        "token is expired",
			tokenString: tokens[1],
			tokenSecret: defaultHmacKey,
			want:        uuids[1],
			wantErr:     true,
		},
		{
			name:        "token secret is invalid ",
			tokenString: tokens[1],
			tokenSecret: "",
			want:        uuids[1],
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
