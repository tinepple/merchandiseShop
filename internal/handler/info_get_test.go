package handler

import (
	"MerchandiseShop/internal/handler/mocks"
	"MerchandiseShop/internal/storage"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestHandler_GetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)

	type expectedResult struct {
		statusCode int
		body       string
	}

	type handlerFields struct {
		storage     Storage
		authService authService
	}

	tests := []struct {
		name          string
		handlerFields handlerFields

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

			expectedResult: expectedResult{
				statusCode: 401,
				body:       `{"errors":"произошла ошибка авторизации"}`,
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

			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка GetPurchasesByUserID",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(1000, nil)
					m.EXPECT().GetPurchasesByUserID(gomock.Any(), 1).Return(make([]storage.Inventory, 0), errors.New("internal error"))
					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)
					m.EXPECT().GetUserID("someToken").Return(1, nil)
					return m
				}(t),
			},

			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка GetTransactionsByUserID",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					inventories := make([]storage.Inventory, 1)
					inventories[0] = storage.Inventory{
						Name:     "test",
						Quantity: 100,
					}
					m := mocks.NewMockStorage(ctrl)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(1000, nil)
					m.EXPECT().GetPurchasesByUserID(gomock.Any(), 1).Return(inventories, nil)
					m.EXPECT().GetTransactionsByUserID(gomock.Any(), 1).Return(make([]storage.CoinsHistory, 0), errors.New("internal error"))
					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)
					m.EXPECT().GetUserID("someToken").Return(1, nil)
					return m
				}(t),
			},

			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Успешно",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					inventories := make([]storage.Inventory, 1)
					inventories[0] = storage.Inventory{
						Name:     "test",
						Quantity: 100,
					}
					transactions := make([]storage.CoinsHistory, 1)
					transactions[0] = storage.CoinsHistory{
						UserNameFrom: "test1",
						UserIDFrom:   1,
						UserNameTo:   "test2",
						UserIDTo:     2,
						Amount:       100,
					}
					m := mocks.NewMockStorage(ctrl)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(1000, nil)
					m.EXPECT().GetPurchasesByUserID(gomock.Any(), 1).Return(inventories, nil)
					m.EXPECT().GetTransactionsByUserID(gomock.Any(), 1).Return(transactions, nil)
					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)
					m.EXPECT().GetUserID("someToken").Return(1, nil)
					return m
				}(t),
			},
			expectedResult: expectedResult{
				statusCode: 200,
				body:       `{"coins":1000,"inventory":[{"type":"test","quantity":100}],"coinHistory":{"received":null,"sent":[{"toUser":"test2","amount":100}]}}`,
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

			req, _ := http.NewRequest("GET", "/api/info", nil)
			req.Header.Set("Authorization", "Bearer someToken")
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedResult.statusCode, w.Code)
			assert.Equal(t, tt.expectedResult.body, w.Body.String())
		})
	}
}
