package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const iterations = 100_000
const keyLen = 32

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type UserRepository struct {
	mu      sync.RWMutex
	byID    map[string]User
	byEmail map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{byID: make(map[string]User), byEmail: make(map[string]string)}
}

func (r *UserRepository) Save(user User) error {
	if user.ID == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("user id, email and password hash are required")
	}
	email := normalizeEmail(user.Email)
	r.mu.Lock()
	defer r.mu.Unlock()
	if existingID, ok := r.byEmail[email]; ok && existingID != user.ID {
		return errors.New("email already registered")
	}
	user.Email = email
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	r.byID[user.ID] = user
	r.byEmail[email] = user.ID
	return nil
}

func (r *UserRepository) FindByEmail(email string) (User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byEmail[normalizeEmail(email)]
	if !ok {
		return User{}, false
	}
	user, ok := r.byID[id]
	return user, ok
}

func (r *UserRepository) FindByID(id string) (User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byID[id]
	return user, ok
}

type Session struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

type SessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{sessions: make(map[string]Session)}
}

func (r *SessionRepository) Save(session Session) error {
	if session.UserID == "" || session.TokenHash == "" {
		return errors.New("session user id and token hash are required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.TokenHash] = session
	return nil
}

func (r *SessionRepository) FindByToken(token string) (Session, bool) {
	hash := HashToken(token)
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[hash]
	if !ok || time.Now().UTC().After(session.ExpiresAt) {
		return Session{}, false
	}
	return session, true
}

type Service struct {
	Users    *UserRepository
	Sessions *SessionRepository
	TTL      time.Duration
}

func NewService() *Service {
	return &Service{Users: NewUserRepository(), Sessions: NewSessionRepository(), TTL: 24 * time.Hour}
}

func (s *Service) Register(id string, email string, password string) (string, User, error) {
	if len(password) < 8 {
		return "", User{}, errors.New("password must be at least 8 characters")
	}
	hash, err := HashPassword(password)
	if err != nil {
		return "", User{}, err
	}
	user := User{ID: id, Email: normalizeEmail(email), PasswordHash: hash}
	if err := s.Users.Save(user); err != nil {
		return "", User{}, err
	}
	token, err := s.createSession(user.ID)
	return token, user, err
}

func (s *Service) Login(email string, password string) (string, User, error) {
	user, ok := s.Users.FindByEmail(email)
	if !ok || !VerifyPassword(password, user.PasswordHash) {
		return "", User{}, errors.New("invalid email or password")
	}
	token, err := s.createSession(user.ID)
	return token, user, err
}

func (s *Service) Authenticate(token string) (User, bool) {
	session, ok := s.Sessions.FindByToken(token)
	if !ok {
		return User{}, false
	}
	return s.Users.FindByID(session.UserID)
}

func (s *Service) createSession(userID string) (string, error) {
	token, err := randomHex(32)
	if err != nil {
		return "", err
	}
	return token, s.Sessions.Save(Session{UserID: userID, TokenHash: HashToken(token), ExpiresAt: time.Now().UTC().Add(s.TTL)})
}

func HashPassword(password string) (string, error) {
	salt, err := randomBytes(16)
	if err != nil {
		return "", err
	}
	key := pbkdf2SHA256([]byte(password), salt, iterations, keyLen)
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iterations, hex.EncodeToString(salt), hex.EncodeToString(key)), nil
}

func VerifyPassword(password string, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2_sha256" {
		return false
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	salt, err := hex.DecodeString(parts[2])
	if err != nil {
		return false
	}
	want, err := hex.DecodeString(parts[3])
	if err != nil {
		return false
	}
	got := pbkdf2SHA256([]byte(password), salt, iter, len(want))
	return hmac.Equal(got, want)
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func randomHex(size int) (string, error) {
	bytes, err := randomBytes(size)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func randomBytes(size int) ([]byte, error) {
	out := make([]byte, size)
	_, err := rand.Read(out)
	return out, err
}

func pbkdf2SHA256(password []byte, salt []byte, iter int, length int) []byte {
	hashLen := 32
	blocks := (length + hashLen - 1) / hashLen
	var out []byte
	for block := 1; block <= blocks; block++ {
		u := pbkdf2Block(password, salt, iter, block)
		out = append(out, u...)
	}
	return out[:length]
}

func pbkdf2Block(password []byte, salt []byte, iter int, block int) []byte {
	mac := hmac.New(sha256.New, password)
	mac.Write(salt)
	mac.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
	u := mac.Sum(nil)
	out := append([]byte(nil), u...)
	for i := 1; i < iter; i++ {
		mac = hmac.New(sha256.New, password)
		mac.Write(u)
		u = mac.Sum(nil)
		for j := range out {
			out[j] ^= u[j]
		}
	}
	return out
}
