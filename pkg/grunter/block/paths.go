package block

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/romainframe/grunter/pkg/env"
	"github.com/romainframe/grunter/pkg/utils"
)

func transformSpecialPath(specialPath string) (string, error) {

	specialPath, err := transformRelativePath(specialPath)
	if err != nil {
		return specialPath, err
	}

	if env.GRUNT_REPO_ROOT != "" {
		specialPath = strings.ReplaceAll(specialPath, env.GRUNT_REPO_ROOT, "${get_repo_root()}")
	}

	return specialPath, nil
}

func transformRelativePath(specialPath string) (string, error) {
	if !strings.HasPrefix(specialPath, "...") {
		return specialPath, nil
	}

	fmt.Printf("specialPath: %s\n", specialPath)

	specialPath = strings.ReplaceAll(specialPath, ".../", "")

	parts := strings.Split(specialPath, "/")
	if len(parts) == 0 {
		return specialPath, fmt.Errorf("invalid special path: %s", specialPath)
	}

	keyword := parts[0]
	remainingParts := parts[1:]

	currentDir, err := os.Getwd()
	if err != nil {
		return specialPath, err
	}

	startDir := filepath.Join(currentDir)

	fmt.Printf("finding upwards from %s with keyword %s\n", startDir, keyword)

	foundDir, err := utils.FindUpwards(startDir, keyword, 20)
	if err != nil {
		return specialPath, err
	}

	return filepath.Join(foundDir, strings.Join(remainingParts, "/")), nil
}
