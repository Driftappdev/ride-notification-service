package fcm

import "time"

type Config struct {
	ProjectID             string
	CredentialsFile       string
	MaxTokensPerBatch     int
	DefaultTTL            time.Duration
	DryRun                bool
	DeleteInvalidTokens   bool
	DefaultAndroidChannel string
}

func (c Config) Normalize() Config {
	if c.MaxTokensPerBatch <= 0 || c.MaxTokensPerBatch > 500 {
		c.MaxTokensPerBatch = 500
	}
	if c.DefaultTTL <= 0 {
		c.DefaultTTL = 24 * time.Hour
	}
	return c
}
