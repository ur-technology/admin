package main

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/homemade/scl"
)

type configurer interface {
	database() databaseConfig
}

type databaseConfig struct {
	Endpoint           string `hcl:"endpoint"`
	ServiceCredentials string `hcl:"service_credentials"`
}

type config struct {
	DatabaseConfig databaseConfig `hcl:"database"`
}

func newConfig(path string, variables ...environmentMap) (configurer, error) {

	c := &config{}

	p, err := scl.NewParser(scl.NewDiskSystem())

	if err != nil {
		return nil, err
	}

	for _, m := range variables {
		for k, v := range m {
			p.SetParam(k, fmt.Sprintf(`%s`, v))
		}
	}

	if err := p.Parse(path); err != nil {
		return nil, err
	}

	if err := hcl.Decode(c, p.String()); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *config) database() databaseConfig {
	return c.DatabaseConfig
}
