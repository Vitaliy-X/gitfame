package main

import (
	"log"

	"gitlab.com/slon/shad-go/gitfame/internal"
	constants "gitlab.com/slon/shad-go/gitfame/internal/util/constants"
	"gitlab.com/slon/shad-go/gitfame/internal/util/load"
)

type ContributorStats struct {
	fileParams *internal.FileParameters
	cmdArgs    *internal.RepoConfig
	mappings   []internal.FileMapping
	statistics *internal.Statistics
}

func NewContributorStats() *ContributorStats {
	cmdArgs := initializeRepoConfig()
	mappings := loadMappings()
	fileParams := retrieveFiles(mappings, cmdArgs)
	return &ContributorStats{fileParams: fileParams, cmdArgs: cmdArgs, mappings: mappings}
}

func initializeRepoConfig() *internal.RepoConfig {
	cmdArgs := internal.InitRepoConfig()
	err := cmdArgs.ConfigureFlags()
	if err != nil {
		log.Fatalf("[ERROR] %s:\n%v", "initializeRepoConfig", err)
	}
	return cmdArgs
}

func loadMappings() []internal.FileMapping {
	return load.LoadMaps()
}

func retrieveFiles(mappings []internal.FileMapping, cmdArgs *internal.RepoConfig) *internal.FileParameters {
	return internal.RetrieveAllFiles(mappings, cmdArgs, cmdArgs.Revision, cmdArgs.Repository)
}

func (cs *ContributorStats) DisplayResults() {
	cs.statistics.SortResults(cs.cmdArgs.SortOrderKey)
	switch cs.cmdArgs.OutputFormat {
	case constants.Tabular:
		cs.statistics.PrintTabular()
	case constants.CSV:
		cs.statistics.PrintCSV()
	case constants.StandardJSON:
		cs.statistics.PrintJSON()
	case constants.JSONLines:
		cs.statistics.PrintJSONLines()
	}
}

func (cs *ContributorStats) CalculateStats() {
	statistics := internal.CalculateStats(cs.fileParams)
	cs.statistics = &statistics
}

func main() {
	stats := NewContributorStats()
	stats.CalculateStats()
	stats.DisplayResults()
}
