package utils

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)
