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

	_ "embed"

	"github.com/Didar1505/project_test.git/internal/auth/providers/otp"
	"github.com/Didar1505/project_test.git/internal/auth/session"
	"github.com/Didar1505/project_test.git/internal/mailer"
	"github.com/Didar1505/project_test.git/internal/user"
)

type Service struct {
	users    user.Repository
	otpRepo  otp.Repository
	sessions session.SessionRepository
	mailer   mailer.Mailer
	tokens   TokenManager

	otpTTL     time.Duration
	accessTTL  time.Duration
	refreshTTL time.Duration
}

//go:embed templates/email_otp.html
var emailOtpTemplate string

func NewService(users user.Repository, otp otp.Repository, sessions session.SessionRepository, mailer mailer.Mailer, tokens TokenManager) *Service {
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

func builtHTMLOTP(code string) string {
	html := strings.ReplaceAll(emailOtpTemplate, "{{CODE}}", code)
	return html
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
	htmlBody := builtHTMLOTP(code)
	textBody := "Your verification code is: " + code
	return s.mailer.SendOTP(ctx, email, htmlBody, textBody)
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

	return s.issueTokens(ctx, u, ua, ip)
}

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

	u, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	access, err := s.tokens.SignAccess(u.ID, s.accessTTL)
	if err != nil {
		return nil, err
	}

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

func (s *Service) LoginWithGoogle(ctx context.Context, email, fullName, subject, ua, ip string) (*AuthResponse, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, errors.New("email required")
	}

	var u *user.User
	var err error

	if subject != "" {
		u, err = s.users.GetByProviderSubject(ctx, "google", subject)
		if err != nil && err != user.ErrNotFound {
			return nil, err
		}
	}

	if u == nil {
		u, err = s.users.GetByEmail(ctx, email)
		if err != nil && err != user.ErrNotFound {
			return nil, err
		}
	}

	if u == nil {
		newUser := &user.User{
			Email:          &email,
			AuthProvider:   "google",
			NativeLanguage: "tk",
		}
		if fullName != "" {
			newUser.FullName = &fullName
		}
		if subject != "" {
			newUser.ProviderSubject = &subject
		}
		if err := s.users.Create(ctx, newUser); err != nil {
			return nil, err
		}
		u = newUser
	}

	return s.issueTokens(ctx, u, ua, ip)
}

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

func (s *Service) issueTokens(ctx context.Context, u *user.User, ua, ip string) (*AuthResponse, error) {
	access, err := s.tokens.SignAccess(u.ID, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshPlain := randomTokenHex(32)
	refreshHash := hashString(refreshPlain)

	uaStr := ua
	ipStr := ip

	sess := &session.Session{
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
