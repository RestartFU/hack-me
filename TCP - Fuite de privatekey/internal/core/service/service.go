package service

import (
	"crypto"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net"
	"tcp/fuitedeprivatekey/internal/config"
	"time"
)

var message = "hi, brenda, this is robert, be careful, hackers are everywhere nowadays, download the joint file to protect your computer!"

var header = `
====================
 PRIVATE KEY LEAK
====================
You should send the message's signature in less than 2 seconds !

note: the hash method used has also leaked, apparently it is MD5

oh no, robert's privatekey leaked: "%s"
make the brenda think that robert sent her the following message: "%s"
`

func Start(cfg config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Network.Host, cfg.Network.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// listener could not accept connection, process is either shutting down
			// or there is a big big problem
			return err
		}

		go func() {
			key, err := rsa.GenerateKey(cryptorand.Reader, 512)
			if err != nil {
				fmt.Printf("could not generate rsa key: %s", err)
				return
			}
			keyBuf := x509.MarshalPKCS1PrivateKey(key)
			b64Key := base64.StdEncoding.EncodeToString(keyBuf)
			_, err = io.WriteString(conn, fmt.Sprintf(header, b64Key, message))
			if err != nil {
				fmt.Printf("could not write to client (%s): %s", conn.RemoteAddr().String(), err)
				return
			}

			// Making sure to set the deadline
			conn.SetDeadline(time.Now().Add(time.Second * 2))

			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					// listener could not accept connection, process is either shutting down,
					// the cdead
					// or there is a big big problem
					return
				}
				buf = buf[:n]

				h := crypto.MD5.New()
				h.Write([]byte(message))
				err = rsa.VerifyPKCS1v15(&key.PublicKey, crypto.MD5, h.Sum(nil), buf)
				if err != nil {
					io.WriteString(conn, "signature does not match message")
					fmt.Printf("could not verify signature (%s): %s", conn.RemoteAddr().String(), err)
					return
				}

				io.WriteString(conn, "congrats! you have infected brenda with your virus, here's your flag: !socialandcyberh#cker#")
			}
		}()
	}
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
