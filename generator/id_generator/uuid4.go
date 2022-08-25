package id_generator

import (
	"github.com/google/uuid"
)

func GetUUIDv4() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
