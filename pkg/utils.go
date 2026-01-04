package pkg

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateCode(prefix string) string {
	prefix = strings.ToUpper(prefix)

	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	timestamp := time.Now().Format("0102")                 // MMDD
	uniquePart := strings.ToUpper(uuid.New().String()[:5]) // first 5 chars of UUID

	return fmt.Sprintf("%s%s%s", prefix, timestamp, uniquePart)
}
