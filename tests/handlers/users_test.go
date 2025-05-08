package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
	apiV1 "todo_list_go/internal/handlers/v1"
	"todo_list_go/internal/service"
	mockService "todo_list_go/internal/service/mocks"
	customErrors "todo_list_go/pkg/errors"
)

func TestUserSignUp(t *testing.T) {
	type mockBehaviour func(s *mockService.MockUser, input service.SignUpUserInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            service.SignUpUserInput
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "Test User", "email": "test@gmail.com", "password": "password123"}`,
			inputUser: service.SignUpUserInput{
				Name:     "Test User",
				Email:    "test@gmail.com",
				Password: "password123",
			},
			mockBehaviour: func(s *mockService.MockUser, input service.SignUpUserInput) {
				s.EXPECT().SignUp(gomock.Any(), input).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: ``,
		},
		{
			name:                 "Wrong data",
			inputBody:            `{}`,
			inputUser:            service.SignUpUserInput{},
			mockBehaviour:        func(s *mockService.MockUser, input service.SignUpUserInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":{"type":"dict","details":{"email":"is required","name":"is required","password":"is required"}}}`,
		},
		{
			name:      "Service error",
			inputBody: `{"name": "Test User", "email": "test@gmail.com", "password": "password123"}`,
			inputUser: service.SignUpUserInput{
				Name:     "Test User",
				Email:    "test@gmail.com",
				Password: "password123",
			},
			mockBehaviour: func(s *mockService.MockUser, input service.SignUpUserInput) {
				s.EXPECT().SignUp(gomock.Any(), input).Return(customErrors.ErrUserAlreadyExists)
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"error":{"type":"string","details":"user with such email already exists"}}`,
		},
		{
			name:      "DB error",
			inputBody: `{"name": "Test User", "email": "test@gmail.com", "password": "password123"}`,
			inputUser: service.SignUpUserInput{
				Name:     "Test User",
				Email:    "test@gmail.com",
				Password: "password123",
			},
			mockBehaviour: func(s *mockService.MockUser, input service.SignUpUserInput) {
				s.EXPECT().SignUp(gomock.Any(), input).Return(errors.New("DB error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":{"type":"string","details":"internal server error"}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			user := mockService.NewMockUser(c)
			testCase.mockBehaviour(user, testCase.inputUser)

			services := &service.Services{Users: user}
			handler := apiV1.NewHandler(services, nil)

			// Init server
			r := gin.New()
			r.POST("api/v1/users/sign-up", handler.SignUp)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/users/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
