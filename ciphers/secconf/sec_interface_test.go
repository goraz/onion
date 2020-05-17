package secconf

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCipher(t *testing.T) {
	Convey("Testing the newcipher", t, func() {
		data := "lorem ipsum"
		pubReader := bytes.NewReader([]byte(pubring))
		b, err := Encode([]byte(data), pubReader)
		So(err, ShouldBeNil)

		secReader := bytes.NewReader([]byte(secring))
		c, err := NewCipher(secReader)
		So(err, ShouldBeNil)
		br, err := c.Decrypt(bytes.NewBuffer(b))
		So(err, ShouldBeNil)
		So(string(br), ShouldEqual, data)
	})
}
