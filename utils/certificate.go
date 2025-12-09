package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCertificateCode → unique cert ID like CERT-20250912-XYZ123
func GenerateCertificateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("CERT-%d-%06d", time.Now().Year(), rand.Intn(1000000))
}
