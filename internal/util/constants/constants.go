package constants

type OutputFormat string

const (
	Tabular      OutputFormat = "tabular"
	CSV          OutputFormat = "csv"
	StandardJSON OutputFormat = "json"
	JSONLines    OutputFormat = "json-lines"
)

type OrderKey string

const (
	Lines   OrderKey = "lines"
	Commits OrderKey = "commits"
	Files   OrderKey = "files"
)
