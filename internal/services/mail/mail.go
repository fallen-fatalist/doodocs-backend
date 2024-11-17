package mail

import "io"

type mailService struct {
}

func NewMailService() *mailService {
	return &mailService{}
}

func (s *mailService) SendFile(file io.Reader, emails []string) error {
	return nil
}
