package GGCaptcha

import "gopkg.in/gomail.v2"

type SMTPClient struct {
	Dialer *gomail.Dialer
	From   string
}

// NewSMTPClient 初始化 SMTP 客户端
func NewSMTPClient(smtpHost string, smtpPort int, username, password string, from string) *SMTPClient {
	return &SMTPClient{
		Dialer: gomail.NewDialer(smtpHost, smtpPort, username, password),
		From:   from,
	}
}

// SendVerificationEmail 发送验证邮件
func (client *SMTPClient) SendVerificationEmail(to, token, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", client.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content) //"Please verify your email: http://yourapp.com/verify?token="+token)

	return client.Dialer.DialAndSend(m)
}
