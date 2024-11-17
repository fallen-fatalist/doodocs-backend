package mail

type mailService struct {
}

func NewMailService() *mailService {
	return &mailService{}
}

func (s *mailService) SendFile(file []byte, emails []string) error {
	return nil
}
