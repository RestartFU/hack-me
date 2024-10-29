package config_test

import (
	"os"
	"tcp/fuitedeprivatekey/internal/config"
	"testing"

	"github.com/restartfu/gophig"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	g := gophig.NewGophig[config.Config]("./../../config.toml", gophig.TOMLMarshaler{}, os.ModePerm)
	_, err := g.LoadConf()
	require.NoError(t, err)
}
