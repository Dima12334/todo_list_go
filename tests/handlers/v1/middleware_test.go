package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
	apiV1 "todo_list_go/internal/handlers/v1"
	mockJwt "todo_list_go/pkg/auth/mocks"
)

func TestUserIdentityMiddleware(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockJwt.MockTokenManager, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockJwt.MockTokenManager, token string) {
				r.EXPECT().ParseJWT(token).Return("9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
		},
		{
			name:                 "Empty Header",
			headerName:           "",
			mockBehavior:         func(r *mockJwt.MockTokenManager, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":{"type":"string","details":"empty auth header"}}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer",
			mockBehavior:         func(r *mockJwt.MockTokenManager, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":{"type":"string","details":"invalid auth header"}}`,
		},
		{
			name:                 "Invalid header",
			headerName:           "Authorization",
			headerValue:          "Beer token",
			token:                "token",
			mockBehavior:         func(r *mockJwt.MockTokenManager, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":{"type":"string","details":"invalid auth header"}}`,
		},
		{
			name:        "Invalid token",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockJwt.MockTokenManager, token string) {
				r.EXPECT().ParseJWT(token).Return("9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d", errors.New("error get user claims from token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":{"type":"string","details":"error get user claims from token"}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			tokenManager := mockJwt.NewMockTokenManager(c)
			testCase.mockBehavior(tokenManager, testCase.token)

			handler := apiV1.NewHandler(nil, tokenManager)

			// Init server
			r := gin.New()
			r.GET("/user-identity", handler.UserIdentityMiddleware, func(c *gin.Context) {
				id, _ := c.Get("userId")
				c.String(200, "%s", id)
			})

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user-identity", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
