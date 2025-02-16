package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHandler_GetTokenFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		args           *gin.Context
		expectedResult string
	}{
		{
			name: "Успешное получение токена",
			args: &gin.Context{
				Request: &http.Request{
					Header: http.Header{
						"Authorization": {"Bearer token"},
					},
				},
			},
			expectedResult: "token",
		},
		{
			name: "Неправильный формат заголовка",
			args: &gin.Context{
				Request: &http.Request{
					Header: http.Header{
						"Authorization": {"Bearertoken"},
					},
				},
			},
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTokenFromRequest(tt.args)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
