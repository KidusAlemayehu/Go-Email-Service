package models

import (
	"email-service/utils/secrets"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID             uuid.UUID             `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name           string                `json:"name" gorm:"not null"`
	Identifier     string                `json:"identifier" gorm:"not null"`
	ImageUrl       string                `json:"image_url"`
	AccessSettings ServiceAccessSettings `json:"access_settings" gorm:"foreignKey:ServiceID"`
	Subscription   ServiceSubscription   `json:"subscription" gorm:"foreignKey:ServiceID"`
	Usage          ServiceUsage          `json:"usage" gorm:"foreignKey:ServiceID"`
	RegisteredAt   time.Time             `json:"registered_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type ServiceAccessSettings struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServiceID   uuid.UUID  `json:"service_id" gorm:"type:uuid;not null"`
	AccessKey   string     `json:"access_key" gorm:"type:text;not null;unique"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"not null"`
	RotatedFrom *uuid.UUID `json:"rotated_from" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s *ServiceAccessSettings) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	if s.AccessKey == "" {
		s.AccessKey = secrets.GenerateSecureKey(32)
	}

	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = time.Now().Add(90 * 24 * time.Hour)
	}

	return
}

func (s *ServiceAccessSettings) RotateKey(tx *gorm.DB) error {
	newKey := &ServiceAccessSettings{
		ServiceID:   s.ServiceID,
		AccessKey:   secrets.GenerateSecureKey(32),
		ExpiresAt:   time.Now().Add(90 * 24 * time.Hour),
		RotatedFrom: &s.ID,
	}

	if err := tx.Create(newKey).Error; err != nil {
		return err
	}

	return nil
}

type ServiceSubscription struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServiceID           uuid.UUID `json:"service_id" gorm:"type:uuid;not null;unique"`
	PlanName            string    `json:"plan_name" gorm:"not null"`
	MonthlyEmailQuota   int       `json:"monthly_email_quota" gorm:"not null"`
	MonthlyStorageQuota int64     `json:"monthly_storage_quota" gorm:"not null"`
	ExpiresAt           time.Time `json:"expires_at"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ServiceUsage struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServiceID   uuid.UUID `json:"service_id" gorm:"type:uuid;not null;unique"`
	EmailsSent  int       `json:"emails_sent" gorm:"default:0"`
	StorageUsed int64     `json:"storage_used" gorm:"default:0"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	LastReset   time.Time `json:"last_reset"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
