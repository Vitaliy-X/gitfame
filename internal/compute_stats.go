package internal

import (
	"os/exec"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/exceptions"
	"gitlab.com/slon/shad-go/gitfame/internal/util/execute"
)

type Statistics struct {
	LinesPerUser       map[string]int
	CommitsPerUser     map[string]map[string]bool
	CommitCountPerUser map[string]int
	FilesPerUser       map[string]map[string]bool
	FileCountPerUser   map[string]int
	AggregatedData     map[string][3]int
	OrderedData        [][4]string
}

var globalStats Statistics

func CalculateStats(fp *FileParameters) Statistics {
	globalStats = Statistics{
		LinesPerUser:       make(map[string]int),
		CommitsPerUser:     make(map[string]map[string]bool),
		CommitCountPerUser: make(map[string]int),
		FilesPerUser:       make(map[string]map[string]bool),
		FileCountPerUser:   make(map[string]int),
		AggregatedData:     make(map[string][3]int),
	}

	for _, filePath := range fp.FileList {
		AnalyzeFile(filePath, fp.Config.Repository, fp.Config.Revision, fp.Config.Commiter)
	}

	for user, commitCount := range globalStats.CommitCountPerUser {
		lineCount := 0
		if actualLineCount, exists := globalStats.LinesPerUser[user]; exists {
			lineCount = actualLineCount
		}
		globalStats.AggregatedData[user] = [3]int{
			lineCount,
			commitCount,
			globalStats.FileCountPerUser[user],
		}
	}

	return globalStats
}

func Register(author string, commitHash string, filePath string) {
	// commit
	if _, ok := globalStats.CommitsPerUser[author]; !ok {
		globalStats.CommitsPerUser[author] = make(map[string]bool)
	}
	if _, ok := globalStats.CommitsPerUser[author][commitHash]; !ok {
		globalStats.CommitsPerUser[author][commitHash] = true
		globalStats.CommitCountPerUser[author]++
	}

	// file
	if _, ok := globalStats.FilesPerUser[author]; !ok {
		globalStats.FilesPerUser[author] = make(map[string]bool)
	}
	if _, ok := globalStats.FilesPerUser[author][filePath]; !ok {
		globalStats.FilesPerUser[author][filePath] = true
		globalStats.FileCountPerUser[author]++
	}
}

func AnalyzeFile(filePath, gitDirectory, commitPointer string, flag bool) {

	commitLog, err := execute.ExecuteGit(exec.Command("git", "blame", "--line-porcelain", "-b", commitPointer, filePath), gitDirectory)

	exceptions.Exception(err, "AnalyzeFile")

	lines := strings.Split(commitLog, "\n")
	var author, commitHash string

	if isFileEmpty(lines) {
		processEmptyFile(commitPointer, filePath, gitDirectory)
	}

	for i := 0; i < len(lines); i++ {
		words := strings.Split(lines[i], " ")

		if flag {
			if words[0] == "committer" {
				commitHash = strings.Split(lines[i-5], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		} else {
			if words[0] == "author" {
				commitHash = strings.Split(lines[i-1], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		}

		globalStats.LinesPerUser[author]++
		Register(author, commitHash, filePath)
	}
}

func processEmptyFile(commitPointer, filePath, gitDirectory string) {

	gitLog, err := execute.ExecuteGit(exec.Command("git", "log", "-p", commitPointer, "--follow", "--", filePath), gitDirectory)

	exceptions.Exception(err, "processEmptyFile")

	logLines := strings.Split(gitLog, "\n")
	commitHash := strings.Split(logLines[0], " ")[1]
	words := strings.Split(logLines[1], " ")
	author := strings.Join(words[1:len(words)-1], " ")

	Register(author, commitHash, filePath)
}

func isFileEmpty(lines []string) bool {
	return len(lines) == 1 && lines[0] == ""
}
