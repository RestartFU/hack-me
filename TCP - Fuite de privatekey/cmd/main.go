package main

import (
	"errors"
	"os"
	"tcp/fuitedeprivatekey/internal/config"
	"tcp/fuitedeprivatekey/internal/core/service"

	"github.com/restartfu/gophig"
)

func main() {
	cfg, err := readConfig("./config.toml")
	if err != nil {
		panic("could not read config: " + err.Error())
	}

	err = service.Start(cfg)
	if err != nil {
		panic("could not read config: " + err.Error())
	}
}

func readConfig(path string) (config.Config, error) {
	g := gophig.NewGophig[config.Config](path, gophig.TOMLMarshaler{}, os.ModePerm)
	cfg, err := g.LoadConf()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg = config.DefaultConfig()
			return cfg, g.SaveConf(cfg)
		}
		return cfg, err
	}

	return cfg, nil
}
