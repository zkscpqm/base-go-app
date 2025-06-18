package common

import (
	"os"
	"strings"
)

func JoinPath(parts ...string) string {
	joined := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		if joined == "" {
			joined = part
		} else {
			joined = joined + string(os.PathSeparator) + part
		}
	}
	return joined
}

func EnsureSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix) + suffix
}
