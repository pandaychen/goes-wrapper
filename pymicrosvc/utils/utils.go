package utils

import (
	"fmt"
)

func GetGrpcResolverScheme(scheme string) string {
	return fmt.Sprintf("%s:///", scheme)
}
