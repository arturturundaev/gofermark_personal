package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) UserExists(login string) (bool, error) {
	args := m.Called(login)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockService) Register(login, password string) (*uuid.UUID, error) {
	args := m.Called(login, password)
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

type MockJWTValidator struct {
	mock.Mock
}

func (m *MockJWTValidator) InitToken(ctx *gin.Context, uid *uuid.UUID) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func TestUserRegisterHandler_Handler(t *testing.T) {
	type bodyStr struct {
		id       string
		login    string
		password string
	}

	type userExistsResponse struct {
		exists bool
		err    error
	}

	type registerResponse struct {
		userID *uuid.UUID
		err    error
	}

	type initTokenResponse struct {
		er error
	}

	userID := uuid.New()
	var tests = []struct {
		name               string
		body               bodyStr
		userExistsResponse userExistsResponse
		register           registerResponse
		initToken          initTokenResponse
		expectedErr        int
	}{
		{
			"Fail get Data From Request Body",
			bodyStr{},
			userExistsResponse{exists: true, err: nil},
			registerResponse{},
			initTokenResponse{},
			http.StatusInternalServerError,
		},
		{
			"Fail when check user in DB",
			bodyStr{login: "UserExists.DB.FailOnCheck", password: "password"},
			userExistsResponse{false, errors.New("FailDB")},
			registerResponse{},
			initTokenResponse{},
			http.StatusInternalServerError,
		},
		{
			"User already exists in DB",
			bodyStr{login: "UserExists", password: "password"},
			userExistsResponse{true, nil},
			registerResponse{},
			initTokenResponse{},
			http.StatusConflict,
		},
		{
			"Fail when save user to DB",
			bodyStr{login: "UserNotExists.DB.FailOnSave", password: "password"},
			userExistsResponse{false, nil},
			registerResponse{nil, errors.New("FailDB")},
			initTokenResponse{},
			http.StatusInternalServerError,
		},
		{
			"Fail init token",
			bodyStr{login: "UserNotExists.FailInitToken", password: "password"},
			userExistsResponse{false, nil},
			registerResponse{&userID, nil},
			initTokenResponse{errors.New("Fail init token")},
			http.StatusInternalServerError,
		},
		{
			"Success",
			bodyStr{login: "UserNotExists.Success", password: "password"},
			userExistsResponse{false, nil},
			registerResponse{&userID, nil},
			initTokenResponse{nil},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := new(MockService)
			jwtValidatorMock := new(MockJWTValidator)
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			ctx.Request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"`+tt.body.login+`","password":"`+tt.body.password+`"}`))

			serviceMock.On("UserExists", tt.body.login).Return(tt.userExistsResponse.exists, tt.userExistsResponse.err)
			serviceMock.On("Register", tt.body.login, tt.body.password).Return(tt.register.userID, tt.register.err)
			jwtValidatorMock.On("InitToken", ctx, tt.register.userID).Return(tt.initToken.er)

			handler := NewUserRegisterHandler(serviceMock, jwtValidatorMock)

			handler.Handler(ctx)

			status := ctx.Writer.Status()
			assert.Equal(t, tt.expectedErr, status)

		})
	}
}
