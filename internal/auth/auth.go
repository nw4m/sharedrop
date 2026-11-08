package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a plain password with a bcrypt hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken makes a random session token
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SessionStore holds active session tokens in memory (simple, no persistence needed)
type SessionStore struct {
	tokens map[string]bool
}

func NewSessionStore() *SessionStore {
	return &SessionStore{tokens: make(map[string]bool)}
}

func (s *SessionStore) Add(token string) {
	s.tokens[token] = true
}

func (s *SessionStore) Valid(token string) bool {
	// constant-time comparison to avoid timing attacks
	for t := range s.tokens {
		if subtle.ConstantTimeCompare([]byte(t), []byte(token)) == 1 {
			return true
		}
	}
	return false
}

func (s *SessionStore) Remove(token string) {
	delete(s.tokens, token)
}

// ParseBasicAuth extracts credentials from Authorization header
func ParseBasicAuth(header string) (string, string, bool) {
	if !strings.HasPrefix(header, "Basic ") {
		return "", "", false
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// FormatSize makes file sizes human-readable
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
