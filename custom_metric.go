package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CustomMetric struct {
	Name       string
	ScriptPath string
}

func NewCustomMetric(script_path string) (*CustomMetric, error) {

	if debug {
		Log(DEBUG_LVL_NOTICE, fmt.Sprintf("Custom metric file: %s", script_path))
	}

	if _, err := os.Stat(script_path); os.IsNotExist(err) {
		return &CustomMetric{}, errors.New(fmt.Sprintf("Script %s does not exists!", script_path))
	}

	_, file := filepath.Split(script_path)
	return &CustomMetric{
		Name:       strings.Split(file, ".")[0],
		ScriptPath: script_path,
	}, nil
}

func (cm *CustomMetric) GetStats() (string, error) {
	out, err := exec.Command(cm.ScriptPath).Output()
	if err != nil {
		fmt.Printf("%s", err)
		return "", err
	}
	return strings.Trim(string(out), "\n"), nil
}
