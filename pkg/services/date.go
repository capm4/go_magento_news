package services

import (
	"time"
)

//return current time in format 2006-01-02T15:04:05
func CurrentDateWithFormat() string {
	return time.Now().Format("2006-01-02T15:04:05")
}
