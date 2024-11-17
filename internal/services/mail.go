package services

type MailService interface {
	SendFile(file []byte, emails []string) error
}
