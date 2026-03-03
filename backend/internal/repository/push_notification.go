package repository

import (
	"elderly-care-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type PushNotificationRepository struct {
	db *gorm.DB
}

func NewPushNotificationRepository(db *gorm.DB) *PushNotificationRepository {
	return &PushNotificationRepository{db: db}
}

// PushToken methods
func (r *PushNotificationRepository) SaveToken(token *model.PushToken) error {
	// 检查是否已存在相同token
	var existing model.PushToken
	err := r.db.Where("token = ? AND user_id = ?", token.Token, token.UserID).First(&existing).Error
	if err == nil {
		// 更新已有记录
		existing.IsActive = true
		existing.LastUsedAt = &[]time.Time{time.Now()}[0]
		return r.db.Save(&existing).Error
	}

	// 创建新记录
	return r.db.Create(token).Error
}

func (r *PushNotificationRepository) GetUserTokens(userID uint, platform string) ([]model.PushToken, error) {
	var tokens []model.PushToken
	query := r.db.Where("user_id = ? AND is_active = ?", userID, true)
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	err := query.Order("last_used_at DESC").Find(&tokens).Error
	return tokens, err
}

func (r *PushNotificationRepository) DeleteToken(token string) error {
	return r.db.Model(&model.PushToken{}).Where("token = ?", token).Update("is_active", false).Error
}

// PushNotification methods
func (r *PushNotificationRepository) CreateNotification(notification *model.PushNotification) error {
	return r.db.Create(notification).Error
}

func (r *PushNotificationRepository) UpdateNotificationStatus(id uint, status string, sentAt *time.Time, failedReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if sentAt != nil {
		updates["sent_at"] = sentAt
	}
	if failedReason != "" {
		updates["failed_reason"] = failedReason
	}
	return r.db.Model(&model.PushNotification{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PushNotificationRepository) GetUserNotifications(userID uint, limit int) ([]model.PushNotification, error) {
	var notifications []model.PushNotification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *PushNotificationRepository) GetUnsentNotifications(limit int) ([]model.PushNotification, error) {
	var notifications []model.PushNotification
	err := r.db.Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}
