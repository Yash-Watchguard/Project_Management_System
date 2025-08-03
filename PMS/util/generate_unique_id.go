package util

import("github.com/google/uuid")

func GenerateUniqueUUID() string{
	return uuid.New().String()
}