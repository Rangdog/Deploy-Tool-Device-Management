package service

import (
	"BE_Manage_device/config"
	"BE_Manage_device/pkg/utils"
	"net/url"
	"strings"
	"sync"

	gomail "gopkg.in/mail.v2"
)

type EmailService struct {
	Dialer *gomail.Dialer
	Mutex  sync.Mutex
}

func NewEmailService(pw string) *EmailService {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "thanhhaxuan02@gmail.com", pw)
	return &EmailService{
		Dialer: dialer,
	}
}

func (service *EmailService) SendActivationEmail(email string, token string, redirectUrl string) error {
	service.Mutex.Lock()
	defer service.Mutex.Unlock()
	m := gomail.NewMessage()
	m.SetHeader("From", "thanhhaxuan02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify Account")
	safeToken := url.QueryEscape(token)
	safeRedirect := url.QueryEscape(redirectUrl)
	actLink := config.BASE_URL_BACKEND + "api/activate?token=" + safeToken + "&redirectUrl=" + safeRedirect
	m.SetBody("text/html", "Click link to activate account: <a href= '"+actLink+"'> Activate </a>")
	for {
		sender, err := service.Dialer.Dial()
		if err != nil {
			return err
		}
		defer sender.Close()
		err = gomail.Send(sender, m)
		if err != nil {
			utils.LogEmailError("activation", email, err)
			if strings.Contains(err.Error(), "broken pipe") {
				continue
			}
		} else {
			utils.LogEmailSuccess("activation", email)
		}
		return err
	}
}

func (service *EmailService) SendEmail(email string, subject string, body string) error {
	service.Mutex.Lock()
	defer service.Mutex.Unlock()
	msg := gomail.NewMessage()
	msg.SetHeader("From", "thanhhaxuan02@gmail.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	for {
		sender, err := service.Dialer.Dial()
		if err != nil {
			return err
		}
		defer sender.Close()
		err = gomail.Send(sender, msg)
		if err != nil {
			utils.LogEmailError("activation", email, err)
			if strings.Contains(err.Error(), "broken pipe") {
				continue
			}
		} else {
			utils.LogEmailSuccess("activation", email)
		}
		return err
	}
}

func (service *EmailService) SendEmails(emails []string, subject string, body string) {
	var wg sync.WaitGroup
	for _, email := range emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			service.SendEmail(email, subject, body)
		}(email)
	}
	wg.Wait()
}
