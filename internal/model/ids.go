package model

import (
	"fmt"
	"strings"
	"time"
)

func MakeID(prefix string, now time.Time, sequence int64) string {
	stamp := now.UTC().Format("20060102-150405.000")
	return fmt.Sprintf("%s-%s-%04d", strings.ToLower(prefix), stamp, sequence)
}

func CleanID(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func RequireID(field, value string) error {
	if CleanID(value) == "" {
		return FieldError{Field: field, Message: "不能为空"}
	}
	return nil
}
