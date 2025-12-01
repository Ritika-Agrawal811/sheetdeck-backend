package entities

type FilePaths struct {
	NewPath string
	OldPath string
}

var Filters = map[string]bool{
	"recent":           true,
	"oldest":           true,
	"most_viewed":      true,
	"least_viewed":     true,
	"most_downloaded":  true,
	"least_downloaded": true,
}
