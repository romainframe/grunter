package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Predefined errors for file operations.
var (
	// ErrFindUpwards is returned when the target directory is not found.
	ErrFindUpwards = fmt.Errorf("target directory not found")
)

// FindUpwards searches for the first part of the path starting from the current directory and moving upwards.
// Returns the path to the found directory or an empty string if not found.
func FindUpwards(startDir, target string, maxDepth int) (string, error) {
	if maxDepth <= 0 {
		return "", fmt.Errorf("folder search max depth reached")
	}

	var foundPath string
	err := filepath.WalkDir(startDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkDir: %w", err)
		}
		if d.IsDir() && d.Name() == target {
			foundPath = path
			return nil
		}
		return nil
	})
	if err != nil {
		return "", WrapError(ErrFindUpwards, err)
	}

	if foundPath != "" {
		return foundPath, nil
	}

	newStartDir := filepath.Dir(startDir)
	if newStartDir == "/" {
		return "", WrapError(ErrFindUpwards, fmt.Errorf("target '%s' not found in '%s'", target, startDir))
	}

	// If the target was not found in the current directory, search downwards
	foundPath, _ = findDownwards(newStartDir, target)
	if foundPath != "" {
		return foundPath, nil
	}

	return FindUpwards(newStartDir, target, maxDepth-1)
}

// findDownwards searches for the remaining parts of the path within the given directory.
func findDownwards(startDir string, target string) (string, error) {
	currentDir := startDir
	parts := strings.Split(target, "/")
	for _, part := range parts {
		found, err := findInDirectory(currentDir, part)
		if err != nil {
			return "", err
		}
		currentDir = found
	}
	return currentDir, nil
}

// findInDirectory searches for a target directory or file within the given directory.
func findInDirectory(dir, target string) (string, error) {
	var targetPath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path) == target {
			targetPath = path
			return nil
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if targetPath != "" {
		return targetPath, nil
	}

	return "", fmt.Errorf("findInDirectory: target '%s' not found in '%s'", target, dir)
}

// FindFileInParentTarget searches a file in a parent target folder.
func FindFileInParentTarget(parentFolder, targetFolder, fileName string, maxDepth int) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Split the parentFolder path into components
	parts := strings.Split(parentFolder, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid parentFolder path")
	}

	// Find the first part of the parentFolder path
	foundPath, err := FindUpwards(currentDir, parts[0], maxDepth)
	if err != nil {
		return "", err
	}

	// Find the remaining parts of the parentFolder path
	if len(parts) > 1 {
		subPath := strings.Join(parts[1:], "/")
		foundPath, err = findDownwards(foundPath, subPath)
		if err != nil {
			return "", err
		}
	}

	// Now find the targetFolder and fileName within the final foundPath
	targetFolderPath, err := findInDirectory(foundPath, targetFolder)
	if err != nil {
		return "", err
	}

	filePath, err := findInDirectory(targetFolderPath, fileName)
	if err != nil {
		return "", err
	}

	return ComputeRelativePath(currentDir, filePath)
}

func FindFileInParent(fileName string, maxDepth int) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath, err := FindUpwards(currentDir, fileName, maxDepth)
	if err != nil {
		return "", err
	}
	return ComputeRelativePath(currentDir, filePath)
}

// ComputeRelativePath computes the relative path from base to target.
// It returns the relative path as a string and any error encountered.
func ComputeRelativePath(basePath, targetPath string) (string, error) {
	// Clean the paths to remove any unnecessary parts
	basePath = filepath.Clean(basePath)
	targetPath = filepath.Clean(targetPath)

	// Compute the relative path
	relPath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return "", err // Return the error to the caller
	}

	return relPath, nil // Return the relative path
}

func DoesFileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
