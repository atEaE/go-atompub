package atompub

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// Authenticator adds authentication to HTTP requests
type Authenticator interface {
	Authenticate(req *http.Request) error
}

// NoAuth implements a no-op authenticator
type NoAuth struct{}

// Authenticate does nothing for NoAuth
func (a *NoAuth) Authenticate(req *http.Request) error {
	// no-op
	return nil
}

var _ Authenticator = (*NoAuth)(nil)

// NewNoAuth creates a new NoAuth authenticator
func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

// WSSEAuth implements WSSE UsernameToken authentication.
// See: https://docs.oasis-open.org/wss-m/wss/v1.1.1/os/wss-UsernameTokenProfile-v1.1.1-os.html
type WSSEAuth struct {
	username string
	password string
}

// Authenticate adds WSSE authentication to the request
func (a *WSSEAuth) Authenticate(req *http.Request) error {
	nonce, err := genNonce()
	if err != nil {
		return fmt.Errorf("generate wsse nonce: %w", err)
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	passwordDigest, err := genPasswordDigest(a.password, nonce, createdAt)
	if err != nil {
		return fmt.Errorf("generate wsse password digest: %w", err)
	}

	req.Header.Set("X-WSSE", genWSSEValue(a.username, passwordDigest, nonce, createdAt))
	req.Header.Set("Authorization", `WSSE profile="UsernameToken"`)
	return nil
}

func genNonce() ([]byte, error) {
	nonce := make([]byte, 20)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func genPasswordDigest(password string, nonce []byte, created string) (string, error) {
	h := sha1.New()
	if _, err := h.Write(nonce); err != nil {
		return "", err
	}
	if _, err := h.Write([]byte(created)); err != nil {
		return "", err
	}
	if _, err := h.Write([]byte(password)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func genWSSEValue(username, passwordDigest string, nonce []byte, created string) string {
	nonceB64 := base64.StdEncoding.EncodeToString(nonce)
	return fmt.Sprintf(
		`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
		username, passwordDigest, nonceB64, created,
	)
}

var _ Authenticator = (*WSSEAuth)(nil)

// NewWSSEAuth creates a new WSSEAuth authenticator
func NewWSSEAuth(username, password string) *WSSEAuth {
	return &WSSEAuth{
		username: username,
		password: password,
	}
}
