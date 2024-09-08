package usecase

import (
	csrf "csrf/ddd/domain"
	"errors"
)

type CsrfTokenRepository interface {
	Save(token *csrf.CsrfToken) error
	Find(sessionID string) (string, error)
}

// CSRFトークン発行ユースケース
type CSRFTokenUseCase struct {
	Repository CsrfTokenRepository
}

func (cc *CSRFTokenUseCase) Generate(sessionID string, secretId string) (*csrf.CsrfToken, error) {
	token, err := csrf.NewCsrfToken(sessionID, []byte(secretId))
	if err != nil {
		return nil, err
	}
	if err := cc.Repository.Save(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (cc *CSRFTokenUseCase) Validate(token string, sessionID string) error {
	savedToken, err := cc.Repository.Find(sessionID)
	if err != nil {
		return err
	}

	if savedToken != token {
		return errors.New("invalid CSRF token")
	}
	return nil
}
