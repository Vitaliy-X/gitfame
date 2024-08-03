package internal

import (
	"os/exec"
	"path/filepath"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/exceptions"
	"gitlab.com/slon/shad-go/gitfame/internal/util/execute"
)

type FileMapping struct {
	Name       string
	FileType   string
	Extensions []string
}

type FileParameters struct {
	FileList []string
	Config   *RepoConfig
	Mappings []FileMapping
}

func RetrieveAllFiles(mappings []FileMapping, config *RepoConfig, commitHash, gitDirectory string) *FileParameters {
	params := &FileParameters{Config: config, Mappings: mappings}

	gitTreeOutput, err := execute.ExecuteGit(exec.Command("git", "ls-tree", "-r", "--name-only", commitHash), gitDirectory)

	exceptions.Exception(err, "retrieveAllFiles")

	fileEntries := strings.Split(gitTreeOutput, "\n")
	for _, entry := range fileEntries {
		if entry != "" {
			if !isLanguageAccepted(entry, params.Mappings, params.Config.Languages, params.Config.Extensions) {
				continue
			}
			if len(params.Config.Exclude) > 0 && matchesAnyPattern(entry, params.Config.Exclude) {
				continue
			}
			if len(params.Config.Restricted) > 0 && !matchesAnyPattern(entry, params.Config.Restricted) {
				continue
			}
			params.FileList = append(params.FileList, entry)
		}
	}
	return params
}

func isLanguageAccepted(filePath string, mappings []FileMapping, acceptedLanguages, allowedExtensions []string) bool {
	fileExtension := filepath.Ext(filePath)
	fileLanguage := ""

	for _, mapping := range mappings {
		for _, ext := range mapping.Extensions {
			if strings.EqualFold(fileExtension, ext) {
				fileLanguage = mapping.Name
				break
			}
		}
		if fileLanguage != "" {
			break
		}
	}

	if len(acceptedLanguages) > 0 {
		if fileLanguage == "" {
			return false
		}
		languageAccepted := false
		for _, language := range acceptedLanguages {
			if strings.EqualFold(language, fileLanguage) {
				languageAccepted = true
				break
			}
		}
		if !languageAccepted {
			return false
		}
	}

	if len(allowedExtensions) > 0 {
		extensionAccepted := false
		for _, ext := range allowedExtensions {
			if strings.EqualFold(fileExtension, ext) {
				extensionAccepted = true
				break
			}
		}
		if !extensionAccepted {
			return false
		}
	}

	return true
}

func matchesAnyPattern(fileName string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := filepath.Match(pattern, fileName)
		if match {
			return true
		}
	}
	return false
}
