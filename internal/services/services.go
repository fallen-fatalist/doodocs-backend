package services

import (
	"zip-api/internal/services/mail"
	"zip-api/internal/services/zipservice"
)

// Global variables
var (
	MailServiceInstance MailService = mail.NewMailService()
	ZipServiceInstance  ZipService  = zipservice.NewZipService()
)
