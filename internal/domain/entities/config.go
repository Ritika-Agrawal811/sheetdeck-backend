package entities

type ResourceLimitConfig struct {
	LimitBytes  int64
	LimitPretty string
}

var ResourceLimits = map[string]ResourceLimitConfig{
	"database": {
		LimitBytes:  500 * 1024 * 1024, // 500 MB
		LimitPretty: "500 MB",
	},
	"storage": {
		LimitBytes:  1 * 1024 * 1024 * 1024, // 1 GB
		LimitPretty: "1 GB",
	},
}

type DatabaseSize struct {
	SizeBytes  int64  `json:"size_bytes"`
	SizePretty string `json:"size_pretty"`
}
