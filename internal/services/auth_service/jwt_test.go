package auth_service

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestHandler_GenerateJWT(t *testing.T) {
	tests := []struct {
		name   string
		userID int
	}{
		{
			name:   "Успешная генераций JWT токена",
			userID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			result, err := s.GenerateJWT(tt.userID)
			assert.NilError(t, err)

			userID, err := s.GetUserID(result)
			assert.NilError(t, err)

			assert.Equal(t, userID, tt.userID)
		})
	}
}
