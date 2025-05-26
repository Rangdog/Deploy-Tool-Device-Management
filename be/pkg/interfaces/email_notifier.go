package interfaces

type EmailNotifier interface {
	SendEmail(to string, subject string, body string) error
	SendEmails(to []string, subject string, body string)
}
