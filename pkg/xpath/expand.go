package xpath

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func Expand(rawPath string) (string, error) {
	switch {
	case strings.HasPrefix(rawPath, "~/"):
		var err error
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to get user home dir: %w", err)
		}
		return path.Join(homeDir, rawPath[2:]), nil
	}
	return rawPath, nil
}
