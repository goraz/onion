package secconf

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/fzerorubigd/onion"
)

type cipher struct {
	secretKeyring []byte
}

func (c *cipher) Decrypt(r io.Reader) ([]byte, error) {
	return Decode(r, bytes.NewReader(c.secretKeyring))
}

func NewCipher(secRing io.Reader) (onion.Cipher, error) {
	b, err := ioutil.ReadAll(secRing)
	if err != nil {
		return nil, err
	}

	return &cipher{
		secretKeyring: b,
	}, nil
}
