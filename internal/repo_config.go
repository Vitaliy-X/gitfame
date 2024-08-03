package internal

import (
	flag "github.com/spf13/pflag"
	"gitlab.com/slon/shad-go/gitfame/internal/util/constants"
)

type RepoConfig struct {
	Repository   string
	Revision     string
	Commiter     bool
	OutputFormat constants.OutputFormat
	SortOrderKey constants.OrderKey
	Languages    []string
	Restricted   []string
	Exclude      []string
	Extensions   []string
}

func InitRepoConfig() *RepoConfig {
	return &RepoConfig{}
}

func (rc *RepoConfig) ConfigureFlags() error {
	var orderBy string
	var outputStyle string

	flag.StringVar(&rc.Repository, "repository", "./", "Path to the repository.")
	flag.StringVar(&rc.Revision, "revision", "HEAD", "Specific commit to use.")
	flag.BoolVar(&rc.Commiter, "use-committer", false, "Consider committer instead of author.")
	flag.StringVar(&outputStyle, "format", "tabular", "Format of the output results.")
	flag.StringVar(&orderBy, "order-by", "lines", "Criteria for sorting results.")
	flag.StringSliceVar(&rc.Languages, "languages", []string{}, "Programming languages to consider.")
	flag.StringSliceVar(&rc.Restricted, "restrict-to", []string{}, "Scope limitations.")
	flag.StringSliceVar(&rc.Exclude, "exclude", []string{}, "Patterns to ignore.")
	flag.StringSliceVar(&rc.Extensions, "extensions", []string{}, "File extensions to include.")

	flag.Parse()

	switch orderBy {
	case "lines":
		rc.SortOrderKey = constants.Lines
	case "commits":
		rc.SortOrderKey = constants.Commits
	case "files":
		rc.SortOrderKey = constants.Files
	}

	switch outputStyle {
	case "tabular":
		rc.OutputFormat = constants.Tabular
	case "csv":
		rc.OutputFormat = constants.CSV
	case "json":
		rc.OutputFormat = constants.StandardJSON
	case "json-lines":
		rc.OutputFormat = constants.JSONLines
	}

	shouldExit, exitValue := validate(rc)
	if shouldExit {
		return exitValue
	}

	return nil
}
