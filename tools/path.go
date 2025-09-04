package tools

import (
	"os"
	"path/filepath"
	"strings"
)

// ensureSafePath 确保 path 在 baseDir 内，防止目录穿越攻击
func EnsureSafePath(baseDir, path string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// 判断是否在 baseDir 内
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", os.ErrPermission
	}
	return absPath, nil
}
