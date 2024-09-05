package csrf

import (
	"context"
	"time"

	"example.com/project/domain/csrf"
	"example.com/project/infra/redis"
)

type CSRFUsecase struct {
	tokenService *csrf.CSRFTokenService
	redisClient  redis.RedisClient
}

func NewCSRFUsecase(tokenService *csrf.CSRFTokenService, redisClient redis.RedisClient) *CSRFUsecase {
	return &CSRFUsecase{
		tokenService: tokenService,
		redisClient:  redisClient,
	}
}

func (u *CSRFUsecase) GenerateToken(ctx context.Context, userID string) (string, error) {
	token, err := u.tokenService.GenerateToken()
	if err != nil {
		return "", err
	}

	// Redisに保存
	err = u.redisClient.Set(ctx, userID, token.Value, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token.Value, nil
}

func (u *CSRFUsecase) ValidateToken(ctx context.Context, userID, providedToken string) error {
	// Redisからトークンを取得
	storedTokenValue, err := u.redisClient.Get(ctx, userID)
	if err != nil {
		return err
	}

	storedToken := csrf.Token{Value: storedTokenValue, ExpiresAt: time.Now().Add(24 * time.Hour)}

	// トークンを検証
	return u.tokenService.ValidateToken(providedToken, storedToken)
}