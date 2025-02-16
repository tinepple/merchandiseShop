package handler

import (
	"MerchandiseShop/internal/handler/mocks"
	"MerchandiseShop/internal/storage"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestHandler_SendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)

	type expectedResult struct {
		statusCode int
		body       string
	}
	type args struct {
		ToUser string
		Amount int
	}
	type handlerFields struct {
		storage     Storage
		authService authService
	}

	tests := []struct {
		name           string
		handlerFields  handlerFields
		args           args
		expectedResult expectedResult
	}{
		{
			name: "Ошибка доступа при проверке токена",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)
					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)
					m.EXPECT().GetUserID("someToken").Return(0, errors.New("internal error"))
					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 401,
				body:       `{"errors":"произошла ошибка авторизации"}`,
			},
		},
		{
			name: "Ошибка валидации amount пустой",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)
					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 0,
			},
			expectedResult: expectedResult{
				statusCode: 400,
				body:       `{"errors":"произошла ошибка валидации: amount is empty"}`,
			},
		},
		{
			name: "Ошибка GetUserBalance",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(0, errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка amount больше баланса",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(100, nil)

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 1000,
			},
			expectedResult: expectedResult{
				statusCode: 400,
				body:       `{"errors":"произошла ошибка валидации: not enough coins on balance"}`,
			},
		},
		{
			name: "Ошибка GetUserByUsername",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(100, nil)

					m.EXPECT().GetUserByUsername(gomock.Any(), "user").Return(storage.User{}, errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка GetUserBalance для UserTo",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(100, nil)

					m.EXPECT().GetUserByUsername(gomock.Any(), "user").Return(storage.User{
						ID:       2,
						Username: "user",
						Password: "pass",
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 2).Return(0, errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка CreateTransaction",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(100, nil)

					m.EXPECT().GetUserByUsername(gomock.Any(), "user").Return(storage.User{
						ID:       2,
						Username: "user",
						Password: "pass",
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 2).Return(100, nil)
					m.EXPECT().CreateTransaction(gomock.Any(), storage.Transaction{
						UserIDFrom: 1,
						UserIDTo:   2,
						Amount:     100,
					}, 0, 200).Return(errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Успешная отправка монет",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(100, nil)

					m.EXPECT().GetUserByUsername(gomock.Any(), "user").Return(storage.User{
						ID:       2,
						Username: "user",
						Password: "pass",
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 2).Return(100, nil)
					m.EXPECT().CreateTransaction(gomock.Any(), storage.Transaction{
						UserIDFrom: 1,
						UserIDTo:   2,
						Amount:     100,
					}, 0, 200).Return(nil)

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				ToUser: "user",
				Amount: 100,
			},
			expectedResult: expectedResult{
				statusCode: 200,
				body:       ``,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(
				tt.handlerFields.storage,
				tt.handlerFields.authService,
			)

			w := httptest.NewRecorder()

			body, _ := json.Marshal(SentCoins{
				ToUser: tt.args.ToUser,
				Amount: tt.args.Amount,
			})
			req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(string(body)))
			req.Header.Set("Authorization", "Bearer someToken")
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedResult.statusCode, w.Code)
			assert.Equal(t, tt.expectedResult.body, w.Body.String())
		})
	}
}
