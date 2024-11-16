package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Global variables
var (
	Port     = "8080"
	Mail     string
	Password string
)

var mailRegex *regexp.Regexp

func Init() error {
	if os.Getenv("PORT") != "" {
		_, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return fmt.Errorf("while converting port into int: %w", err)
		} else {
			Port = os.Getenv("PORT")
		}
	}
	mailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	Mail = os.Getenv("MAIL")
	if Mail != "" && !mailRegex.Match([]byte(Mail)) {
		return fmt.Errorf("mail does not match the standard format")
	}

	Password = os.Getenv("PASSWORD")

	return nil
}
