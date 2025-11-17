package biz

import (
	"context"
	"fmt"
	"net/smtp"
     "notificationservice/internal/conf"
	"time"
)

type Notification struct {
	
	BookingID uint64
	Email     string
	Subject   string
	Body      string
	Status    string
	CreatedAt time.Time
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

type NotificationRepo interface {
	SaveNotification(ctx context.Context, notif *Notification) error
}

type NotificationUsecase struct {
	repo       NotificationRepo
	emailConf  *EmailConfig
}

func NewNotificationUsecase(repo NotificationRepo, emailConf *EmailConfig) *NotificationUsecase {
	return &NotificationUsecase{
		repo:      repo,
		emailConf: emailConf,
	}
}

func ProvideEmailConfig(c *conf.Email) *EmailConfig {
       if c == nil {
        panic("email config is nil") // helpful guard
    }
    return &EmailConfig{
        SMTPHost: c.SmtpHost,
        SMTPPort: int(c.SmtpPort),
        Username: c.Username,
        Password: c.Password,
        From:     c.From,
    }
}

// Send saves notification to DB and sends real email
func (uc *NotificationUsecase) Send(ctx context.Context, notif *Notification) error {
	notif.CreatedAt = time.Now()

	// 1️⃣ Save to DB
	if err := uc.repo.SaveNotification(ctx, notif); err != nil {
		return err
	}

	// 2️⃣ Send email
	auth := smtp.PlainAuth(
		"",
		uc.emailConf.Username,
		uc.emailConf.Password,
		uc.emailConf.SMTPHost,
	)

	to := []string{notif.Email}
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", notif.Subject, notif.Body))

	addr := fmt.Sprintf("%s:%d", uc.emailConf.SMTPHost, uc.emailConf.SMTPPort)
	if err := smtp.SendMail(addr, auth, uc.emailConf.From, to, msg); err != nil {
		return err
	}

	return nil
}
