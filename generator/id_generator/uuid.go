package id_generator

import (
	"encoding/hex"

	"github.com/rogpeppe/fastuuid"
)

var generator *fastuuid.Generator

func initID() (err error) {
	generator, err = fastuuid.NewGenerator()
	return
}

func NextID() string {
	id := generator.Next()
	return hex.EncodeToString(id[:])
}

func init() {
	initID()
}
