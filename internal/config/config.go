package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// Config holds all the information required by promgen to generate metrics.
type Config struct {
}

func GetDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, "promgen.yaml")
}
