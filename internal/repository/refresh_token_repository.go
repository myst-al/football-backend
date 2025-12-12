package repository

import (
	"football-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Save(userID uint, token string, jti string, expiresAt time.Time) error
	Get(token string) (*models.RefreshToken, error)
	Delete(token string) error
	DeleteByUser(userID uint) error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepo{db}
}

func (r *refreshTokenRepo) Save(userID uint, token, jti string, expiresAt time.Time) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		JTI:       jti,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(&rt).Error
}

func (r *refreshTokenRepo) Get(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepo) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepo) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
