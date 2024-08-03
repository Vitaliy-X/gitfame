package internal

import (
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/slon/shad-go/gitfame/internal/util/constants"
)

func doesCommitExist(commitHash, gitDirectory string) bool {
	cmd := exec.Command("git", "show", commitHash)
	cmd.Dir = gitDirectory
	return cmd.Run() == nil
}

func filePathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)

}
func isOrderKeyValid(orderKey constants.OrderKey) bool {
	validOrderKeys := []constants.OrderKey{constants.Lines, constants.Commits, constants.Files}
	for _, key := range validOrderKeys {
		if orderKey == key {
			return true
		}
	}
	return false
}

func isOutputFormatValid(outputFormat constants.OutputFormat) bool {
	validFormats := []constants.OutputFormat{constants.Tabular, constants.CSV, constants.StandardJSON, constants.JSONLines}
	for _, format := range validFormats {
		if outputFormat == format {
			return true
		}
	}
	return false
}

func validate(config *RepoConfig) (bool, error) {
	if !doesCommitExist(config.Revision, config.Repository) {
		return true, fmt.Errorf("the commit '%s' does not exist in the repository '%s'", config.Revision, config.Repository)
	}
	if !filePathExists(config.Repository) {
		return true, fmt.Errorf("the repository path '%s' is invalid or does not exist", config.Repository)
	}
	if !isOrderKeyValid(config.SortOrderKey) {
		return true, fmt.Errorf("the sort order key '%s' is invalid. valid options are: 'lines', 'commits', 'files'", config.SortOrderKey)
	}
	if !isOutputFormatValid(config.OutputFormat) {
		return true, fmt.Errorf("the output format '%s' is invalid. valid options are: 'tabular', 'csv', 'json', 'json-lines'", config.OutputFormat)
	}
	return false, nil
}
