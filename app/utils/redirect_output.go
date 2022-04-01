package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func RedirectOutput(logPath string) (err error) {
	dir := filepath.Dir(logPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if err != nil {
		return
	}
	err = syscall.Dup2(int(logFile.Fd()), syscall.Stdout)
	if err != nil {
		return
	}
	err = syscall.Dup2(int(logFile.Fd()), syscall.Stderr)
	if err != nil {
		return
	}
	return
}

func RedirectPath(dir, prefix, suffix string) (directpath string) {
	t := time.Now()
	filename := fmt.Sprintf("%s%d%02d%02d-%02d%02d%02d%s", prefix, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), suffix)
	return filepath.Join(dir, filename)
}
