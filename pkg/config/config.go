package config

import (
	"github.com/ekristen/libnuke/pkg/config"
)

// New creates a new extended configuration from a file. This is necessary because we are extended the default
// libnuke configuration to contain additional attributes that are specific to the AWS Nuke tool.
func New(opts config.Options) (*Config, error) {
	// Step 1 - Create the libnuke config
	cfg, err := config.New(opts)
	if err != nil {
		return nil, err
	}

	// Step 2 - Instantiate the extended config
	c := &Config{}

	// Step 3 - Load the same config file against the extended config
	if err := c.Load(opts.Path); err != nil {
		return nil, err
	}

	// Step 4 - Set the libnuke config on the extended config
	c.Config = *cfg

	return c, nil
}

// Config is an extended configuration implementation that adds some additional features on top of the libnuke config.
type Config struct {
	// Config is the underlying libnuke configuration.
	config.Config `yaml:",inline"`
}
