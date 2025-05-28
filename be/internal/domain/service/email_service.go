package service

import (
	"BE_Manage_device/config"
	"BE_Manage_device/pkg/utils"
	"fmt"
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
	const maxRetry = 3
	for attempt := 1; attempt <= maxRetry; attempt++ {
		m := gomail.NewMessage()
		m.SetHeader("From", "thanhhaxuan02@gmail.com")
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Verify Account")
		safeToken := url.QueryEscape(token)
		safeRedirect := url.QueryEscape(redirectUrl)
		actLink := config.BASE_URL_BACKEND + "api/activate?token=" + safeToken + "&redirectUrl=" + safeRedirect
		m.SetBody("text/html", "Click link to activate account: <a href= '"+actLink+"'> Activate </a>")
		service.Mutex.Lock()
		sender, err := service.Dialer.Dial()
		service.Mutex.Unlock()
		if err != nil {
			if attempt == maxRetry {
				return err
			}
			continue
		}
		err = gomail.Send(sender, m)
		sender.Close()
		if err == nil {
			utils.LogEmailSuccess("email", email)
			return nil
		}
		utils.LogEmailError("email", email, err)
		if strings.Contains(err.Error(), "broken pipe") && attempt < maxRetry {
			continue
		} else {
			return err
		}
	}
	return fmt.Errorf("send email retry failed")
}

func (service *EmailService) SendEmail(email string, subject string, body string) error {
	const maxRetry = 3
	for attempt := 1; attempt <= maxRetry; attempt++ {
		msg := gomail.NewMessage()
		msg.SetHeader("From", "thanhhaxuan02@gmail.com")
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)
		service.Mutex.Lock()
		sender, err := service.Dialer.Dial()
		service.Mutex.Unlock()
		if err != nil {
			if attempt == maxRetry {
				return err
			}
			continue
		}
		err = gomail.Send(sender, msg)
		sender.Close()
		if err == nil {
			utils.LogEmailSuccess("email", email)
			return nil
		}
		utils.LogEmailError("email", email, err)
		if strings.Contains(err.Error(), "broken pipe") && attempt < maxRetry {
			continue
		} else {
			return err
		}
	}
	return fmt.Errorf("send email retry failed")
}

func (service *EmailService) SendEmails(emails []string, subject string, body string) {
	const workerCount = 10
	type emailJob struct {
		Email string
	}
	jobs := make(chan emailJob, len(emails))
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				_ = service.SendEmail(job.Email, subject, body)
			}
		}()
	}
	for _, email := range emails {
		jobs <- emailJob{Email: email}
	}
	close(jobs)
	wg.Wait()
}
