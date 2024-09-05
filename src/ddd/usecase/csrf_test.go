package application

import (
	"testing"
	"time"
	"yourproject/infrastructure"

	"github.com/stretchr/testify/mock"
)

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) SaveToken(token string, expiration time.Duration) error {
	args := m.Called(token, expiration)
	return args.Error(0)
}

func (m *MockTokenStorage) GetToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestGenerateToken(t *testing.T) {
	mockStorage := new(MockTokenStorage)
	mockStorage.On("SaveToken", mock.Anything, mock.Anything).Return(nil)

	service := NewCSRFService(mockStorage)
	token, err := service.GenerateToken()

	if err != nil || token == nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateToken(t *testing.T) {
	mockStorage := new(MockTokenStorage)
	mockStorage.On("GetToken", mock.Anything).Return("valid_token", nil)

	service := NewCSRFService(mockStorage)
	err := service.ValidateToken("valid_token")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
