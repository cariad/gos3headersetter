package gos3headersetter

// When describes a header value to set for a given filename extension.
type When struct {
	Extension string `yaml:"extension"`
	Then      string `yaml:"then"`
}
