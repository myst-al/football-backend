package service

import (
	"strings"
	"time"

	"football-backend/internal/config"
	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password, role string) error
	Login(username, password string) (string, string, *models.User, error)
	GetProfile(userID uint) (*models.User, error)
	Refresh(refreshToken string) (string, string, *models.User, error)
	Logout(refreshToken string) error
}

type authService struct {
	repo   repository.UserRepository
	rtRepo repository.RefreshTokenRepository
}

func NewAuthService(userRepo repository.UserRepository, rtRepo repository.RefreshTokenRepository) AuthService {
	return &authService{repo: userRepo, rtRepo: rtRepo}
}

func (s *authService) Register(username, password, role string) error {
	role = strings.ToUpper(role)

	validRoles := map[string]bool{
		"ADMIN":  true,
		"STAFF":  true,
		"VIEWER": true,
	}

	if !validRoles[role] {
		return apperror.NewValidationError("role tidak valid (ADMIN, STAFF, VIEWER)")
	}

	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return apperror.NewConflictError("username sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternalError("gagal proses password")
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := s.repo.Create(user); err != nil {
		return apperror.NewInternalError("gagal membuat user")
	}
	return nil
}

func (s *authService) Login(username, password string) (string, string, *models.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", "", nil, apperror.NewNotFoundError("user tidak ditemukan")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", nil, apperror.NewValidationError("password salah")
	}

	access, refresh, expiresAt, jti, err := generateTokens(user)
	if err != nil {
		return "", "", nil, apperror.NewInternalError("gagal membuat token")
	}

	_ = s.rtRepo.DeleteByUser(user.ID)
	if err := s.rtRepo.Save(user.ID, refresh, jti, expiresAt); err != nil {
		return "", "", nil, apperror.NewInternalError("gagal menyimpan refresh token")
	}

	return access, refresh, user, nil
}

func (s *authService) GetProfile(id uint) (*models.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("user tidak ditemukan")
	}
	return u, nil
}

func generateTokens(user *models.User) (string, string, time.Time, string, error) {
	now := time.Now()
	accessExp := now.Add(15 * time.Minute)
	refreshExp := now.Add(7 * 24 * time.Hour)
	jti := uuid.NewString()

	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"ver":     user.TokenVersion,
		"exp":     accessExp.Unix(),
		"iat":     now.Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"type":    "refresh",
		"exp":     refreshExp.Unix(),
		"iat":     now.Unix(),
		"jti":     jti,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	atSigned, err := at.SignedString([]byte(config.Load().JWTSecret))
	if err != nil {
		return "", "", time.Time{}, "", err
	}

	rtSigned, err := rt.SignedString([]byte(config.Load().JWTSecret))
	if err != nil {
		return "", "", time.Time{}, "", err
	}

	return atSigned, rtSigned, refreshExp, jti, nil
}

func (s *authService) Refresh(refreshToken string) (string, string, *models.User, error) {
	rt, err := s.rtRepo.Get(refreshToken)
	if err != nil {
		return "", "", nil, apperror.NewUnauthorizedError("refresh token tidak dikenal")
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = s.rtRepo.Delete(refreshToken)
		return "", "", nil, apperror.NewUnauthorizedError("refresh token sudah kadaluarsa")
	}

	user, err := s.repo.GetByID(rt.UserID)
	if err != nil {
		return "", "", nil, apperror.NewNotFoundError("user tidak ditemukan")
	}

	user.TokenVersion++
	if err := s.repo.Update(user); err != nil {
		return "", "", nil, apperror.NewInternalError("gagal update token version")
	}

	access, newRefresh, expiresAt, newJti, err := generateTokens(user)
	if err != nil {
		return "", "", nil, apperror.NewInternalError("gagal membuat token baru")
	}

	_ = s.rtRepo.Delete(refreshToken)
	if err := s.rtRepo.Save(user.ID, newRefresh, newJti, expiresAt); err != nil {
		return "", "", nil, apperror.NewInternalError("gagal menyimpan refresh token baru")
	}

	return access, newRefresh, user, nil
}

func (s *authService) Logout(refreshToken string) error {
	rt, err := s.rtRepo.Get(refreshToken)
	if err != nil {
		return apperror.NewUnauthorizedError("refresh token tidak valid")
	}

	user, err := s.repo.GetByID(rt.UserID)
	if err != nil {
		return apperror.NewNotFoundError("user tidak ditemukan")
	}

	user.TokenVersion++
	if err := s.repo.Update(user); err != nil {
		return apperror.NewInternalError("gagal update token version")
	}

	if err := s.rtRepo.Delete(refreshToken); err != nil {
		return apperror.NewInternalError("gagal menghapus refresh token")
	}

	return nil
}
