package circleci

import (
	"github.com/hashicorp/vault/sdk/framework"
	"strings"
)

// Config is the stored configuration.
type Config struct {
	APIToken string   `json:"api-token"`
	OrgId string   `json:"org-id"`
}

// DefaultConfig returns a config with the default values.
func DefaultConfig() *Config {
	return &Config{
		APIToken: "",
		OrgId: "",
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

	if v, ok := d.GetOk("org-id"); ok {
		nv := strings.TrimSpace(v.(string))
		if nv != c.OrgId {
			c.OrgId = nv
			changed = true
		}
	}


	return changed, nil
}
