package egobee

import (
	"encoding/json"
	"errors"
	"regexp"
	"sync"
	"time"
)

var (
	// ErrInvalidDuration is returned from UnmarshalJSON when the JSON does not
	// represent a Duration.
	ErrInvalidDuration = errors.New("invalid duration")

	hasUnitRx = regexp.MustCompile("[a-zA-Z]+")
)

// Scope of a token.
type Scope string

// Possible Scopes.
// See https://www.ecobee.com/home/developer/api/documentation/v1/auth/auth-intro.shtml
var (
	ScopeSmartRead  Scope = "smartRead"
	ScopeSmartWrite Scope = "smartWrite"
	ScopeEMSWrite   Scope = "ems"
)

// TokenDuration wraps time.Duration to add JSON (un)marshalling
type TokenDuration struct {
	time.Duration
}

// MarshalJSON returns JSON representation of Duration.
func (d TokenDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// UnmarshalJSON returns a Duration from JSON representation. Since the ecobee
// API returns durations in Seconds, values will be treated as such.
func (d *TokenDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Second * time.Duration(value)
	case string:
		if !hasUnitRx.Match([]byte(value)) {
			value = value + "s"
		}
		dv, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		d.Duration = dv
	default:
		return ErrInvalidDuration
	}
	return nil
}

// TokenRefreshResponse is returned by the ecobee API on toke refresh.
// See https://www.ecobee.com/home/developer/api/documentation/v1/auth/token-refresh.shtml
type TokenRefreshResponse struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    TokenDuration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	Scope        Scope         `json:"scope"`
}

// TokenStore for ecobee Access and Refresh tokens.
type TokenStore interface {
	// AccessToken gets the access token from the store.
	AccessToken() string

	// RefreshToken gets the refresh token from the store.
	RefreshToken() string

	// ValidFor reports how much longer the access token is valid.
	ValidFor() time.Duration
}

// memoryStore implements tokenStore backed only by memory.
type memoryStore struct {
	mu           sync.RWMutex // protects the following members
	accessToken  string
	refreshToken string
	validUntil   time.Time
}

func (s *memoryStore) AccessToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.accessToken
}

func (s *memoryStore) RefreshToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.refreshToken
}

func (s *memoryStore) ValidFor() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().Sub(s.validUntil)
}

// NewMemoryTokenStore is a TokenStore with no persistence.
func NewMemoryTokenStore(r *TokenRefreshResponse) TokenStore {
	return &memoryStore{
		accessToken:  r.AccessToken,
		refreshToken: r.RefreshToken,
		validUntil:   time.Now().Add(r.ExpiresIn.Duration),
	}
}
