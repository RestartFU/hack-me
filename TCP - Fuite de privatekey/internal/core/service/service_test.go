package service_test

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"tcp/fuitedeprivatekey/internal/config"
	"tcp/fuitedeprivatekey/internal/core/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	t.Run("the service starts and a client connects to it", func(t *testing.T) {
		cfg := config.DefaultConfig()
		go service.Start(cfg)

		// wait for the service to start
		<-time.After(time.Second / 2)

		addr := fmt.Sprintf("%s:%d", cfg.Network.Host, cfg.Network.Port)
		conn, err := net.Dial("tcp", addr)
		require.NoError(t, err)
		var read bool

		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			require.NoError(t, err)

			buf = buf[:n]
			require.Equal(t, len(buf), n)
			log(t, string(buf))

			if read {
				return
			}
			read = true

			regex := regexp.MustCompile(`"(.*?)"`)
			matches := regex.FindAllString(string(buf), 2)
			require.Len(t, matches, 2)

			m1, err := strconv.Unquote(matches[0])
			log(t, m1)
			require.NoError(t, err)
			privKey, err := base64.StdEncoding.DecodeString(m1)
			require.NoError(t, err)

			m2, err := strconv.Unquote(matches[1])
			log(t, m2)

			unmarshaledKey, err := x509.ParsePKCS1PrivateKey([]byte(privKey))
			require.NoError(t, err)
			_ = unmarshaledKey

			h := crypto.MD5.New()
			h.Write([]byte(m2))
			signature, err := unmarshaledKey.Sign(nil, h.Sum(nil), crypto.MD5)
			require.NoError(t, err)

			log(t, string(signature))
			conn.Write(signature)
		}
	})
}

var logIndex int

func log(t *testing.T, s string) {
	f := fmt.Sprintf("\nLOG-%d:\n%s\n\n", logIndex, s)
	_, err := io.WriteString(os.Stdout, f)
	require.NoError(t, err)

	logIndex++
}
