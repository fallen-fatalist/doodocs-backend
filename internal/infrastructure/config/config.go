package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Global variables
var (
	Port                   = "8080"
	BodyLimitInBytes int64 = 1024<<20 + 1024 // 1GB
	Mail             string
	Password         string
)

var mailRegex *regexp.Regexp

func Init() error {
	// Port handling
	if os.Getenv("PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return fmt.Errorf("while converting port into int: %w", err)
		} else {
			if port < 1024 || port > 65535 {
				return fmt.Errorf("invalid port number: must lie in range between 1023 and 65536")
			}
			Port = os.Getenv("PORT")
		}
	}

	// Mail handling
	mailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	Mail = os.Getenv("EMAIL")
	if Mail != "" && !mailRegex.Match([]byte(Mail)) {
		return fmt.Errorf("mail does not match the standard format")
	}

	// Password
	Password = os.Getenv("PASSWORD")

	// Request body limit
	if os.Getenv("BODYLIMIT") != "" {
		bodyLimit, err := strconv.Atoi(os.Getenv("BODYLIMIT"))
		if err != nil {
			return fmt.Errorf("while converting body limit into int: %w", err)
		} else if bodyLimit <= 0 {
			return fmt.Errorf("body limit cannot be negative or zero")
		} else {
			BodyLimitInBytes = int64(bodyLimit)
		}
	}
	return nil
}
