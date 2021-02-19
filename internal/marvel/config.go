package marvel

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ApiBaseUrl    string `envconfig:"MARVEL_API_BASE_URL"`
	ApiKeyPublic  string `envconfig:"MARVEL_API_KEY_PUBLIC"`
	ApiKeyPrivate string `envconfig:"MARVEL_API_KEY_PRIVATE"`

	EagerLoadCache bool `envconfig:"EAGER_LOAD_CACHE"`
}

func NewConfig() *Config {
	c := new(Config)
	envconfig.MustProcess("", c)
	c.validate()
	return c
}

func (c *Config) validate() {
	required := []string{
		c.ApiBaseUrl,
		c.ApiKeyPublic,
		c.ApiKeyPrivate,
	}

	for _, v := range required {
		if v == "" {
			panic("environment variables not properly configured")
		}
	}
}
