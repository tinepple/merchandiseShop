package handler

import (
	"MerchandiseShop/internal/handler/mocks"
	"MerchandiseShop/internal/storage"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestHandler_ItemBuy(t *testing.T) {
	ctrl := gomock.NewController(t)

	type expectedResult struct {
		statusCode int
		body       string
	}
	type args struct {
		item string
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
				item: "test",
			},
			expectedResult: expectedResult{
				statusCode: 401,
				body:       `{"errors":"произошла ошибка авторизации"}`,
			},
		},
		{
			name: "Ошибка GetItem",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetItem(gomock.Any(), "pen").Return(storage.Item{}, errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				item: "pen",
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка GetUserBalance",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetItem(gomock.Any(), "test").Return(storage.Item{
						ID:    1,
						Name:  "test",
						Price: 100,
					}, nil)
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
				item: "test",
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Ошибка userBalance < item.Price",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetItem(gomock.Any(), "test").Return(storage.Item{
						ID:    1,
						Name:  "test",
						Price: 100,
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(10, nil)

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				item: "test",
			},
			expectedResult: expectedResult{
				statusCode: 400,
				body:       `{"errors":"произошла ошибка валидации: not enough coins on balance"}`,
			},
		},
		{
			name: "Ошибка CreatePurchase",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetItem(gomock.Any(), "test").Return(storage.Item{
						ID:    1,
						Name:  "test",
						Price: 100,
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(200, nil)
					m.EXPECT().CreatePurchase(gomock.Any(), 1, 1, 100).Return(errors.New("internal error"))

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				item: "test",
			},
			expectedResult: expectedResult{
				statusCode: 500,
				body:       `{"errors":"произошла внутрення ошибка"}`,
			},
		},
		{
			name: "Успешная покупка",
			handlerFields: handlerFields{
				storage: func(t *testing.T) Storage {
					m := mocks.NewMockStorage(ctrl)

					m.EXPECT().GetItem(gomock.Any(), "test").Return(storage.Item{
						ID:    1,
						Name:  "test",
						Price: 100,
					}, nil)
					m.EXPECT().GetUserBalance(gomock.Any(), 1).Return(200, nil)
					m.EXPECT().CreatePurchase(gomock.Any(), 1, 1, 100).Return(nil)

					return m
				}(t),
				authService: func(t *testing.T) authService {
					m := mocks.NewMockauthService(ctrl)

					m.EXPECT().GetUserID("someToken").Return(1, nil)

					return m
				}(t),
			},
			args: args{
				item: "test",
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

			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/buy/%s", tt.args.item), nil)
			req.Header.Set("Authorization", "Bearer someToken")
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedResult.statusCode, w.Code)
			assert.Equal(t, tt.expectedResult.body, w.Body.String())
		})
	}
}
