package gos3headersetter

// Rule describes a rule for setting a header value.
type Rule struct {
	Header string `yaml:"header"`
	When   []When `yaml:"when"`
	Else   string `yaml:"else"`
}
