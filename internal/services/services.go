package services

import (
	"zip-api/internal/services/mail"
	"zip-api/internal/services/zip"
)

// Global variables
var (
	MailServiceInstance MailService = mail.NewMailService()
	ZipServiceInstance  ZipService  = zip.NewZipService()
)
