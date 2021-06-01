package uuid

import "github.com/google/uuid"

func CreateSlug() string {
	return uuid.New().String()
}

func IsCreatedSlug(value string) bool {
	n := len(value)

	if n > 36 || n < 32 {
		return false
	}

	_, err := uuid.Parse(value)

	return err == nil
}
