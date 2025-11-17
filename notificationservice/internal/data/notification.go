package data

import (
	"context"
	"notificationservice/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type Notification struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	
	BookingID uint64
	Email     string
	Subject   string
	Body      string `gorm:"type:text"`
	Status    string
}

// notificationRepo implements biz.NotificationRepo
type notificationRepo struct {
	data *Data
	log  *log.Helper
}

// NewNotificationRepo returns a biz.NotificationRepo backed by GORM
func NewNotificationRepo(data *Data, logger log.Logger) biz.NotificationRepo {
	return &notificationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveNotification persists a notification into DB
func (r *notificationRepo) SaveNotification(ctx context.Context, notif *biz.Notification) error {
	dbNotif := &Notification{
		
		BookingID: notif.BookingID,
		Email:     notif.Email,
		Subject:   notif.Subject,
		Body:      notif.Body,
		Status:    notif.Status,
	}
	return r.data.db.WithContext(ctx).Create(dbNotif).Error
}
