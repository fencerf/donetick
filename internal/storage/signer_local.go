package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"

	"donetick.com/core/config"
)

type URLSignerLocal struct {
	Secret []byte
}

func NewURLSignerLocal(config *config.Config) *URLSignerLocal {
	return &URLSignerLocal{Secret: []byte(config.Jwt.Secret)}
}

// sign method without expiration:
func (s *URLSignerLocal) Sign(rawPath string) (string, error) {
	sig := s.sign(rawPath)
	values := url.Values{}
	values.Set("sig", sig)
	return fmt.Sprintf("%s?%s", rawPath, values.Encode()), nil
}

func (s *URLSignerLocal) sign(path string) string {
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write([]byte(path))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *URLSignerLocal) IsValid(rawPath string, providedSig string) bool {

	expectedSig := s.sign(rawPath)
	return hmac.Equal([]byte(expectedSig), []byte(providedSig))
}
