package usecase_test

import (
	csrf "csrf/ddd/domain"
	"csrf/ddd/usecase"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCsrfTokenRepositoryは、リポジトリのモック
type MockCsrfTokenRepository struct {
	mock.Mock
}

func (m *MockCsrfTokenRepository) Save(token *csrf.CsrfToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockCsrfTokenRepository) Find(sessionID string) (string, error) {
	args := m.Called(sessionID)
	return args.String(0), args.Error(1)
}

// テストケースを表す構造体
type generateTokenTestCase struct {
	name          string
	sessionID     string
	secretKey     string
	mockSetup     func(*MockCsrfTokenRepository)
	expectedErr   error
	expectedToken *csrf.CsrfToken
}

type validateTokenTestCase struct {
	name        string
	sessionID   string
	token       string
	mockSetup   func(*MockCsrfTokenRepository)
	expectedErr error
}

func TestGenerateToken(t *testing.T) {
	testCases := []generateTokenTestCase{
		{
			name:      "Success",
			sessionID: "test-session",
			secretKey: "secret-key",
			mockSetup: func(repo *MockCsrfTokenRepository) {
				repo.On("Save", mock.Anything).Return(nil)
			},
			expectedErr: nil,
			expectedToken: &csrf.CsrfToken{ // モックデータに合わせて修正
				Token:     "generated-token",
				SessionID: "test-session",
				CreatedAt: time.Now(),
			},
		},
		{
			name:      "Save Failure",
			sessionID: "test-session",
			secretKey: "secret-key",
			mockSetup: func(repo *MockCsrfTokenRepository) {
				repo.On("Save", mock.Anything).Return(errors.New("failed to save token"))
			},
			expectedErr:   errors.New("failed to save token"),
			expectedToken: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockCsrfTokenRepository)
			tc.mockSetup(mockRepo)

			useCase := &usecase.CSRFTokenUseCase{
				Repository: mockRepo,
			}

			token, err := useCase.Generate(tc.sessionID, tc.secretKey)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				log.Printf("!!!!!!!!!!!!!!!!!!!!\n %#v \n !!!!!!!!!!!!!!!", token)
				assert.Equal(t, tc.expectedToken, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestValidateToken(t *testing.T) {
	testCases := []validateTokenTestCase{
		{
			name:      "Validation Success",
			sessionID: "test-session",
			token:     "valid-token",
			mockSetup: func(repo *MockCsrfTokenRepository) {
				repo.On("Find", "test-session").Return("valid-token", nil)
			},
			expectedErr: nil,
		},
		{
			name:      "Validation Failure",
			sessionID: "test-session",
			token:     "invalid-token",
			mockSetup: func(repo *MockCsrfTokenRepository) {
				repo.On("Find", "test-session").Return("valid-token", nil)
			},
			expectedErr: errors.New("invalid CSRF token"),
		},
		{
			name:      "Find Failure",
			sessionID: "test-session",
			token:     "any-token",
			mockSetup: func(repo *MockCsrfTokenRepository) {
				repo.On("Find", "test-session").Return("", errors.New("repository error"))
			},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockCsrfTokenRepository)
			tc.mockSetup(mockRepo)

			useCase := &usecase.CSRFTokenUseCase{
				Repository: mockRepo,
			}

			err := useCase.Validate(tc.token, tc.sessionID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
