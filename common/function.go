package common

import (
	"log"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	log.Println("validating id")
	_, err := uuid.Parse(u)
	return err == nil
}
