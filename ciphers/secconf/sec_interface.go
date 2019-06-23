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
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Decode(data, bytes.NewReader(c.secretKeyring))
}

// NewCipher create a new cipher based on the secconf encoding as specified in the following
// format:
//   base64(gpg(gzip(data)))
func NewCipher(secRing io.Reader) (onion.Cipher, error) {
	b, err := ioutil.ReadAll(secRing)
	if err != nil {
		return nil, err
	}

	return &cipher{
		secretKeyring: b,
	}, nil
}
