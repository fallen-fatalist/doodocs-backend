package services

import "io"

type MailService interface {
	SendFile(file io.Reader, emails []string) error
}
