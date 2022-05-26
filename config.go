package circleci

import (
	"github.com/hashicorp/vault/sdk/framework"
	"strings"
)

// Config is the stored configuration.
type Config struct {
	APIToken string   `json:"api-token"`
}

// DefaultConfig returns a config with the default values.
func DefaultConfig() *Config {
	return &Config{
		APIToken: "",
	}
}

// Update updates the configuration from the given field data.
func (c *Config) Update(d *framework.FieldData) (bool, error) {
	if d == nil {
		return false, nil
	}

	changed := false

	if v, ok := d.GetOk("api-token"); ok {
		nv := strings.TrimSpace(v.(string))
		if nv != c.APIToken {
			c.APIToken = nv
			changed = true
		}
	}

	return changed, nil
}
