package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
	"github.com/Didar1505/project_test.git/internal/user"
)

type Service struct {
	users    user.Repository
	otpRepo  OTPRepository
	sessions SessionRepository
	mailer   Mailer
	tokens   TokenManager

	otpTTL     time.Duration
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(users user.Repository, otp OTPRepository, sessions SessionRepository, mailer Mailer, tokens TokenManager) *Service {
	return &Service{
		users:      users,
		otpRepo:    otp,
		sessions:   sessions,
		mailer:     mailer,
		tokens:     tokens,
		otpTTL:     10 * time.Minute,
		accessTTL:  20 * time.Minute,
		refreshTTL: 30 * 24 * time.Hour,
	}
}

func (s *Service) RequestOTP(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	if email == "" {
		return errors.New("email required")
	}

	code := generate6Digits()
	codeHash := hashString(code)

	if err := s.otpRepo.Create(ctx, email, codeHash, time.Now().UTC().Add(s.otpTTL)); err != nil {
		return err
	}

	return s.mailer.SendOTP(ctx, email, code)
}

func (s *Service) VerifyOTP(ctx context.Context, email, code, ua, ip string) (*AuthResponse, error) {
	email = normalizeEmail(email)

	row, err := s.otpRepo.GetLatestValid(ctx, email, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if row.AttemptsLeft <= 0 {
		return nil, errors.New("otp_attempts_exceeded")
	}
	if hashString(code) != row.CodeHash {
		_ = s.otpRepo.DecrementAttempts(ctx, row.ID)
		return nil, errors.New("invalid_otp")
	}
	_ = s.otpRepo.DeleteByID(ctx, row.ID)

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if err != user.ErrNotFound {
			return nil, err
		}
		newUser := &user.User{
			Email:          &email,
			AuthProvider:   "email_otp",
			NativeLanguage: "tk",
		}
		if err := s.users.Create(ctx, newUser); err != nil {
			return nil, err
		}
		u = newUser
	}

	access, err := s.tokens.SignAccess(u.ID, s.accessTTL)
	if err != nil {
		return nil, err
	}

	// refresh token: сохраняем хэш в sessions
	refreshPlain := randomTokenHex(32)
	refreshHash := hashString(refreshPlain)

	uaStr := ua
	ipStr := ip

	sess := &Session{
		UserID:           u.ID,
		RefreshTokenHash: refreshHash,
		UserAgent:        &uaStr,
		IP:               &ipStr,
		ExpiresAt:        time.Now().UTC().Add(s.refreshTTL),
	}
	if err := s.sessions.Create(ctx, sess); err != nil {
		return nil, err
	}

	_ = s.users.UpdateLastLogin(ctx, u.ID)

	return &AuthResponse{
		AccessToken:  access,
		RefreshToken: refreshPlain,
		User:         user.UserToResponse(*u),
	}, nil
}

// Refresh — проверяем refresh, ротируем, выдаём новый access+refresh
func (s *Service) Refresh(ctx context.Context, refreshToken, ua, ip string) (*AuthResponse, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidRefresh
	}

	now := time.Now().UTC()
	refreshHash := hashString(refreshToken)

	sess, err := s.sessions.GetByRefreshHash(ctx, refreshHash, now)
	if err != nil {
		return nil, err
	}

	// Load user
	u, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	access, err := s.tokens.SignAccess(u.ID, s.accessTTL)
	if err != nil {
		return nil, err
	}

	// Rotation: генерим новый refresh и обновляем запись сессии
	newRefreshPlain := randomTokenHex(32)
	newRefreshHash := hashString(newRefreshPlain)
	newExpires := now.Add(s.refreshTTL)

	if err := s.sessions.Rotate(ctx, sess.ID, newRefreshHash, newExpires, now); err != nil {
		return nil, err
	}

	_ = s.users.UpdateLastLogin(ctx, u.ID)

	_ = ua
	_ = ip

	return &AuthResponse{
		AccessToken:  access,
		RefreshToken: newRefreshPlain,
		User:         user.UserToResponse(*u),
	}, nil
}


// Logout — отзываем текущий refresh
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidRefresh
	}

	now := time.Now().UTC()
	refreshHash := hashString(refreshToken)

	sess, err := s.sessions.GetByRefreshHash(ctx, refreshHash, now)
	if err != nil {
		return err
	}

	return s.sessions.Revoke(ctx, sess.ID, now)
}

// helpers

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func generate6Digits() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	n = n % 1000000
	return leftPad6(n)
}

func leftPad6(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

func randomTokenHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}