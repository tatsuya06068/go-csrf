package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type CsrfToken struct {
	Token     string
	SessionID string
	CreatedAt time.Time
}

// 新しいCSRFトークンを生成
func NewCsrfToken(sessionID string, secretKey []byte) (*CsrfToken, error) {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(sessionID))
	token := hex.EncodeToString(h.Sum(nil))

	return &CsrfToken{
		Token:     token,
		SessionID: sessionID,
		CreatedAt: time.Now(),
	}, nil
}
