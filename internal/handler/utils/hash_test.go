package utils

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHandler_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Успешное хэширование пароля",
			password: "test-password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashPassword(tt.password)
			assert.NilError(t, err)
			assert.Equal(t, CheckPasswordHash(tt.password, result), true)
		})
	}
}
